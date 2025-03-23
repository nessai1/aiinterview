-- +goose Up
-- +goose StatementBegin
ALTER TABLE interview ADD COLUMN feedback TEXT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE interview DROP COLUMN IF EXISTS feedback;
-- +goose StatementEnd
