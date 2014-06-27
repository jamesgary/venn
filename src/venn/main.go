package main

import (
  "database/sql"
  "fmt"
  "github.com/codegangsta/martini"
  _ "github.com/lib/pq"
  "log"
  "net/http"
  "strconv"
  "strings"
)

func main() {
  db, err := sql.Open("postgres", "dbname=venn sslmode=disable")
  if err != nil {
    log.Println("Error opening postgres connection!")
    log.Fatal(err)
  }

  server := martini.Classic()

  server.Get("/movies/:term", func(params martini.Params, w http.ResponseWriter) string {
    w.Header().Set("Content-Type", "application/json")

    rows, err := db.Query(
      "SELECT title FROM movies WHERE title % $1 ORDER BY char_length(tags) DESC LIMIT 10;",
      params["term"],
    )
    if err != nil {
      log.Println("Error querying postgres!")
      log.Fatal(err)
    }

    var movieList []string
    for rows.Next() {
      var title string
      if err := rows.Scan(&title); err != nil {
        log.Fatal(err)
      }
      movieList = append(movieList, strconv.Quote(title))
    }
    return fmt.Sprintf("[%s]", strings.Join(movieList, ","))
  })

  server.Get("/tags/:title", func(params martini.Params, w http.ResponseWriter) string {
    w.Header().Set("Content-Type", "application/json")

    rows, err := db.Query(
      "SELECT tags FROM movies WHERE title = $1",
      params["title"],
    )
    if err != nil {
      log.Println("Error querying postgres!")
      log.Fatal(err)
    }

    var allTagsString string
    for rows.Next() {
      if err := rows.Scan(&allTagsString); err != nil {
        log.Fatal(err)
      }
    }
    tagList := strings.Split(allTagsString, " ")
    for i, tag := range tagList {
      tagList[i] = strconv.Quote(tag)
    }
    return fmt.Sprintf("[%s]", strings.Join(tagList, ","))
  })

  server.Run()
}
