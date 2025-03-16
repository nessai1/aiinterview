-- +goose Up
-- +goose StatementBegin
ALTER TABLE question ADD COLUMN done BOOLEAN NOT NULL DEFAULT false;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE question DROP COLUMN IF EXISTS done;
-- +goose StatementEnd
