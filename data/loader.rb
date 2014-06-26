# assumes movie.list is a list of movies with tags in this format:
# Batman: The Dark Knight (2008)     badass

def insert(conn, title, tags)
  res = conn.exec_params(
    "INSERT INTO movies (title, tags) VALUES ($1, $2)",
    [title, tags.join(" ")]
  )
end

require 'pg'

conn = PG::Connection.open(:dbname => 'venn', :port => 5432 )

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
