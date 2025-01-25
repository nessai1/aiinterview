-- +goose Up
-- +goose StatementBegin
-- Что-бы проверять на коллизии
CREATE TABLE IF NOT EXISTS users (
    uuid UUID NOT NULL PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS interview (
    uuid UUID NOT NULL PRIMARY KEY,
    owner UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    topics JSON NOT NULL -- не самое крутое решение, но для MVP я думаю сойдет
);

CREATE INDEX IF NOT EXISTS ix_interview_owner ON interview (owner);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS ix_interview_owner;
DROP TABLE IF EXISTS interview;
DROP TABLE IF EXISTS users;

-- +goose StatementEnd
