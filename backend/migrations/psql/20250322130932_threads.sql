-- +goose Up
-- +goose StatementBegin
ALTER TABLE interview ADD COLUMN thread VARCHAR(255) NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE interview DROP COLUMN IF EXISTS thread;
-- +goose StatementEnd
