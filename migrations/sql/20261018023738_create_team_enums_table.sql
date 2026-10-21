-- +goose Up
CREATE TABLE IF NOT EXISTS team_enums (
    id BIGSERIAL PRIMARY KEY,
    team_id BIGINT REFERENCES teams(id),
    enum VARCHAR(255) NOT NULL,
    data VARCHAR(255)[] NOT NULL DEFAULT '{}',
    default_datum VARCHAR(255) NULL DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_team_enums UNIQUE (team_id, enum)
);

-- +goose Down
DROP TABLE IF EXISTS team_enums;
