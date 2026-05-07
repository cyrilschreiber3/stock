-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_timestamps()
RETURNS trigger AS $$
BEGIN
    -- Keep created_at immutable
    NEW.created_at := OLD.created_at;

    -- Always bump updated_at when an UPDATE happens
    NEW.updated_at := NOW();

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS set_timestamps();
-- +goose StatementEnd
