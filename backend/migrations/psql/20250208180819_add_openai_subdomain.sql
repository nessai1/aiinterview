-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS assistants (
    id VARCHAR(255) NOT NULL PRIMARY KEY,
    external_id VARCHAR(255) NOT NULL,
    model VARCHAR(80) NOT NULL
);

ALTER TABLE interview ADD COLUMN thread_id VARCHAR(255);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS assistants;
ALTER TABLE interview DROP COLUMN IF EXISTS thread_id;
-- +goose StatementEnd
