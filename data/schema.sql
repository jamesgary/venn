CREATE TABLE movies (title varchar(255), tags text);
CREATE EXTENSION pg_trgm;
CREATE INDEX movie_trgm_idx ON movies USING gist (t gist_trgm_ops);
