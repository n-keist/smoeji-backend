-- +goose Up
-- +goose StatementBegin
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE "public"."users" (
		"id" uuid NOT NULL DEFAULT gen_random_uuid(),
		"email" text NOT NULL,
		"password" text NOT NULL,
		"updated_at" timestamptz NOT NULL,
		"created_at" timestamptz NOT NULL,
		PRIMARY KEY ("id")
	);
CREATE TABLE "public"."refresh_tokens" (
		"id" text NOT NULL DEFAULT gen_random_uuid(),
		"user_id" uuid NOT NULL,
		"value" text NOT NULL,
		"created_at" timestamptz,
		"expires_at" timestamptz NOT NULL DEFAULT now(),
		CONSTRAINT "fk_users_refresh_tokens" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id"),
		PRIMARY KEY ("id")
	);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
