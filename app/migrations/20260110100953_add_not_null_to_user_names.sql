-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ALTER COLUMN first_name SET NOT NULL;
ALTER TABLE users ALTER COLUMN last_name SET NOT NULL;



-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE users ALTER COLUMN first_name DROP NOT NULL;
ALTER TABLE users ALTER COLUMN last_name DROP NOT NULL;



-- +goose StatementEnd
