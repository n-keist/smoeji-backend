-- +goose Up
-- +goose StatementBegin

ALTER TABLE users REPLICA IDENTITY FULL;
ALTER TABLE refresh_tokens REPLICA IDENTITY FULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
