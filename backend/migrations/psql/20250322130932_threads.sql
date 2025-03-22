-- +goose Up
-- +goose StatementBegin
ALTER TABLE interview ADD COLUMN thread NOT NULL VARCHAR(255) DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE interview DROP COLUMN IF EXISTS thread;
-- +goose StatementEnd
