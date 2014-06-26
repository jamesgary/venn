CREATE TABLE movies (title varchar(255), tags text);
CREATE INDEX ON movies ((lower(title)));
CREATE EXTENSION pg_trgm;
