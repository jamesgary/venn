CREATE TABLE movies (id bigserial PRIMARY KEY, title varchar(255) UNIQUE, plot text);
CREATE TABLE keywords (id bigserial PRIMARY KEY, keyword varchar(255) UNIQUE);
CREATE TABLE movie_keywords (movie_id bigint, keyword_id bigint);

CREATE INDEX mk_movie_idx ON movie_keywords(movie_id);
CREATE INDEX mk_keyword_idx ON movie_keywords(keyword_id);
