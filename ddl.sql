CREATE TABLE IF NOT EXISTS pessoas (
    id UUID PRIMARY KEY,
    apelido varchar(32) UNIQUE NOT NULL,
    nome varchar(100) NOT NULL,
    nascimento DATE NOT NULL,
    stack text[] NULL
);
