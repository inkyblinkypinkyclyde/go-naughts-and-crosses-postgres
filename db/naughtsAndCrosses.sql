DROP TABLE IF EXISTS grid;

DROP TABLE IF EXISTS turns;

CREATE TABLE grid (
    id SERIAL PRIMARY KEY,
    value VARCHAR(255)
);

CREATE TABLE turns (
    id SERIAL PRIMARY KEY,
    value VARCHAR(255)
);

