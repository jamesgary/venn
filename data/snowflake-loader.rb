# assumes movie.list is a list of movies with tags in this format:
# Batman: The Dark Knight (2008)     badass

def insert(conn, title, keywords)
  conn.exec_params(
    "INSERT INTO movies
        (title)
      SELECT $1
      WHERE
        NOT EXISTS (
          SELECT title FROM movies WHERE title = $2
        );",
    [title, title]
  )
  movie_id = nil
  conn.exec_params(
    "SELECT id FROM movies where title = $1",
    [title]
  ) do |result|
    result.each do |row|
      movie_id = row.values_at('id')[0]
    end
  end

  keywords.each do |k|
    conn.exec_params(
      "INSERT INTO keywords
          (keyword)
        SELECT $1
        WHERE
          NOT EXISTS (
            SELECT keyword FROM keywords WHERE keyword = $2
          );",
      [k, k]
    )
    keyword_id = nil
    conn.exec_params(
      "SELECT id FROM keywords where keyword = $1",
      [k]
    ) do |result|
      result.each do |row|
        keyword_id = row.values_at('id')[0]
      end
    end
    conn.exec_params(
      "INSERT INTO movie_keywords (movie_id, keyword_id) VALUES ($1, $2)",
      [movie_id.to_i, keyword_id.to_i]
    )
  end
end

require 'pg'

conn = PG::Connection.open(:dbname => 'snowflake', :port => 5432 )

File.open("movies.list", "r") do |f|
  title, prev_title = nil
  tags = []
  f.each_line do |line|
    line = line.encode('UTF-8', 'binary', invalid: :replace, undef: :replace, replace: '')
    words = line.split(" ")
    tag = words.pop
    title = words.join(" ")

    prev_title = title unless prev_title # ran once at first
    if title == prev_title
      tags << tag
    else
      insert(conn, prev_title, tags) # insert for previous movie
      tags = [tag]
      prev_title = title
    end
  end

  # insert the last one
  insert(conn, title, tags)
end
