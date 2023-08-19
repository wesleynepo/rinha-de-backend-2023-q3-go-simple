package data

import (
	"context"
	"database/sql"
	"time"
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



