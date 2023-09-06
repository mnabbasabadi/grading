-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TYPE scale_type AS ENUM ('default','4.0','4.3','5.0','7.0','10.0','ECTS');

CREATE TABLE scale
(
    id SERIAL PRIMARY KEY,
    min FLOAT NOT NULL,
    type scale_type NOT NULL,
    gpa varchar(10) NOT NULL
);

CREATE INDEX scale_type_idx ON scale (type);

INSERT INTO scale (min, gpa,type) VALUES (0.0, 'F', 'default');
INSERT INTO scale (min, gpa,type) VALUES (1.0, 'D', 'default');
INSERT INTO scale (min, gpa,type) VALUES (2.0, 'C', 'default');
INSERT INTO scale (min, gpa,type) VALUES (3.0, 'B', 'default');
INSERT INTO scale (min, gpa,type) VALUES (4.0, 'A', 'default');

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP INDEX IF EXISTS scale_type_idx;
DROP TABLE IF EXISTS scale cascade;
DROP TYPE IF EXISTS scale_type;
