-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE grade
(
    id                    SERIAL PRIMARY KEY,
    student_id            UUID                NOT NULL,
    course_id             UUID                NOT NULL,
    grade                 INTEGER             NOT NULL,
    created_at            TIMESTAMP   NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMP   NOT NULL DEFAULT NOW()
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE IF EXISTS grade cascade;