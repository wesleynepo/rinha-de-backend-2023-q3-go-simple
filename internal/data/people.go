package data

import (
	"context"
	"database/sql"
	"errors"
	"gopherinha/internal/validator"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

var (
    ErrDuplicateApelido = errors.New("duplicate apelido")
)

type PeopleModel struct {
    DB *sql.DB
}

func (p PeopleModel) Count() (int, error) {
    query := `SELECT count(*) FROM pessoas`

    ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
    defer cancel()

    var count int

    err := p.DB.QueryRowContext(ctx, query).Scan(&count)

    if err != nil {
        return 0, err
    }

    return count, nil
}

func (p PeopleModel) Get(id string) (*Person, error) {
    if id == "" {
        return nil, ErrRecordNotFound
    }
    
    query := `
    SELECT id, nome, apelido, nascimento, stack 
    FROM pessoas
    WHERE id = $1`

    var person Person

    ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
    defer cancel()

    err := p.DB.QueryRowContext(ctx, query, id).Scan(
        &person.UUID,
        &person.Nome,
        &person.Apelido,
        &person.Nascimento,
        pq.Array(&person.Stack),
    )

    if err != nil {
        switch {
        case errors.Is(err, sql.ErrNoRows):
            return nil, ErrRecordNotFound
        default:
            return nil, err
        }
    }

    return &person, nil
}

func (p PeopleModel) Insert(person *Person) error {
    query := `
    INSERT INTO pessoas (id, nome, apelido, nascimento, stack) 
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id`

    args := []interface{}{uuid.New(), person.Nome, person.Apelido, person.Nascimento, pq.Array(person.Stack)}

    ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
    defer cancel()

    err := p.DB.QueryRowContext(ctx, query, args...).Scan(&person.UUID)

    if err != nil {
        switch {
        case err.Error() == `pq: duplicate key value violates unique constraint "pessoas_apelido_key"`:
            return ErrDuplicateApelido
        default: 
            return err
        }
    }

    return nil
}

func (p PeopleModel) GetAll(searchTerm string) ([]*Person, error) {
    query := `
    SELECT id, nome, apelido, nascimento, stack 
    FROM pessoas
    WHERE (to_tsvector('simple', nome) @@ plainto_tsquery('simple', $1) OR 
    (to_tsvector('simple', apelido) @@ plainto_tsquery('simple', $1)) OR
    (to_tsvector('simple', array_to_string(stack, ',')) @@ plainto_tsquery('simple', $1)))
    LIMIT 50`

    ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
    defer cancel()

    rows, err := p.DB.QueryContext(ctx, query, searchTerm)
    if err != nil {
        return nil, err
    }

    defer rows.Close()
    people := []*Person{}

    for rows.Next() {
        var person Person 

        err := rows.Scan(
            &person.UUID,
            &person.Nome,
            &person.Apelido,
            &person.Nascimento,
            pq.Array(&person.Stack),
        )

        if err != nil {
            return nil, err
        }

        people = append(people, &person)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return people, nil
}


type Person struct {
    UUID string `json:"id"`
    Apelido string `json:"apelido"`
    Nome string `json:"nome"`
    Nascimento string `json:"nascimento"`
    Stack []string `json:"stack"`
}

func ValidatePerson(v *validator.Validator, person *Person) {
    v.Check(person.Nome != "", "nome", "must be provided")
    v.Check(len(person.Nome) <= 100, "nome", "must not be more than 100 bytes long")

    v.Check(person.Apelido != "", "apelido", "must be provided")
    v.Check(len(person.Apelido) <= 32, "apelido", "must not be more than 32 bytes long")

    v.Check(person.Nascimento != "", "nascimento", "must be provided")
    v.Check(validator.Birthdate(person.Nascimento), "nascimento", "must be YYYY-MM-DD format")
}
