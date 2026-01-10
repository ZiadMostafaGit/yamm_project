-- +goose Up
-- +goose StatementBegin

ALTER TABLE users ALTER COLUMN age SET NOT NULL;
ALTER TABLE users ALTER COLUMN created_at SET NOT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users ALTER COLUMN age DROP NOT NULL;
ALTER TABLE users ALTER COLUMN created_at DROP NOT NULL;
-- +goose StatementEnd
