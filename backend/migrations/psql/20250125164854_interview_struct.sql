-- +goose Up
-- +goose StatementBegin
-- Что-бы проверять на коллизии
CREATE TABLE IF NOT EXISTS users (
    uuid UUID NOT NULL PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS interview (
    uuid UUID NOT NULL PRIMARY KEY,
    owner_uuid UUID NOT NULL REFERENCES users(uuid) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL
);

CREATE INDEX IF NOT EXISTS ix_interview_owner ON interview (owner);

CREATE TABLE IF NOT EXISTS section (
    uuid UUID NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    position INT NOT NULL,
    interview_uuid UUID NOT NULL,
    color VARCHAR(8) NOT NULL,
    is_started BOOLEAN NOT NULL DEFAULT false,
    is_complete BOOLEAN NOT NULL DEFAULT false,
    mark TEXT NULL
);

CREATE INDEX IF NOT EXISTS ix_section_interview ON section (interview_uuid);

CREATE TABLE IF NOT EXISTS question (
    uuid UUID NOT NULL PRIMARY KEY,
    section_uuid UUID NOT NULL PRIMARY KEY,
    position INT NOT NULL,
    question TEXT NOT NULL,
    answer TEXT NULL,
    mark TEXT NULL
);

CREATE INDEX IF NOT EXISTS ix_question_section ON question (section_uuid);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS ix_interview_owner;
DROP INDEX IF EXISTS ix_section_interview;
DROP INDEX IF EXISTS ix_question_section;
DROP TABLE IF EXISTS question;
DROP TABLE IF EXISTS section;
DROP TABLE IF EXISTS interview;
DROP TABLE IF EXISTS users;

-- +goose StatementEnd
