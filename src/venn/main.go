package main

import (
  "database/sql"
  "fmt"
  "github.com/codegangsta/martini"
  _ "github.com/lib/pq"
  "log"
  "net/http"
  "os"
  "strconv"
  "strings"
)

func main() {
  psqlHost := os.Getenv("VENN_PSQL_HOST")
  assetsDir := os.Getenv("VENN_ASSETS_DIR")

  psqlConfig := "dbname=venn sslmode=disable"
  if psqlHost != "" {
    psqlConfig = fmt.Sprintf("host=%s %s", psqlHost, psqlConfig)
  }
  db, err := sql.Open("postgres", psqlConfig)
  if err != nil {
    log.Println("Error opening postgres connection!")
    log.Fatal(err)
  }

  server := martini.Classic()

  if assetsDir != "" {
    server.Use(martini.Static(assetsDir))
  }
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
