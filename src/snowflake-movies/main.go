package main

import (
  "database/sql"
  "fmt"
  "github.com/codegangsta/martini"
  _ "github.com/lib/pq"
  "log"
  "net/http"
  "os"
  "sort"
  "strconv"
  "strings"
)

func main() {
  psqlHost := os.Getenv("VENN_PSQL_HOST")
  assetsDir := os.Getenv("VENN_ASSETS_DIR")

  psqlConfig := "dbname=snowflake sslmode=disable"
  if psqlHost != "" {
    psqlConfig = fmt.Sprintf("host=%s %s", psqlHost, psqlConfig)
  }
  db, err := sql.Open("postgres", psqlConfig)
  if err != nil {
    log.Println("Error opening postgres connection!")
    log.Fatal(err)
  }
  err = db.Ping()
  if err != nil {
    log.Println("Error pinging!")
    log.Fatal(err)
  }

  server := martini.Classic()

  if assetsDir != "" {
    server.Use(martini.Static(assetsDir))
  }
  server.Get("/new_game", func(params martini.Params, w http.ResponseWriter) string {
    w.Header().Set("Content-Type", "application/json")

    var (
      numMoviesWithKeyword int
      keywordId            int
      keyword              string
      movieIds             []int
    )

    for numMoviesWithKeyword < 5 {
      keywordId, keyword = getRandomKeyword(db)
      movieIds = getMovieIdsWithKeywordId(keywordId, db)
      numMoviesWithKeyword = len(movieIds)
    }

    possibleKeywords := getKeywordsForMovies(movieIds, db)
    sort.Strings(possibleKeywords)

    for i, keyword := range possibleKeywords {
      possibleKeywords[i] = strconv.Quote(keyword)
    }

    return fmt.Sprintf(`{"keyword": %s, "movieCount": %d, "possibleKeywords": [%s]}`, strconv.Quote(keyword), len(movieIds), strings.Join(possibleKeywords, ","))
  })

  server.Run()
}

func getRandomKeyword(db *sql.DB) (int, string) {
  rows, err := db.Query("select id, keyword from keywords order by random() limit 1")
  if err != nil {
    log.Println("Error querying postgres!")
    log.Fatal(err)
  }

  var keyword string
  var keywordId int
  for rows.Next() {
    if err := rows.Scan(&keywordId, &keyword); err != nil {
      log.Fatal(err)
    }
  }
  return keywordId, keyword
}

func getMovieIdsWithKeywordId(keywordId int, db *sql.DB) []int {
  rows, err := db.Query("select movie_id from movie_keywords where keyword_id = $1", keywordId)
  if err != nil {
    log.Println("Error querying postgres!")
    log.Fatal(err)
  }

  var movieIds []int
  for rows.Next() {
    var movieId int
    if err := rows.Scan(&movieId); err != nil {
      log.Fatal(err)
    }
    movieIds = append(movieIds, movieId)
  }
  return movieIds
}

func getKeywordsForMovies(movieIds []int, db *sql.DB) []string {
  args := make([]interface{}, len(movieIds))
  qp := make([]string, len(movieIds))
  for i, v := range movieIds {
    args[i] = interface{}(v)
    qp[i] = fmt.Sprintf("$%d", i+1)
  }

  rows, err := db.Query(`
    select keywords.keyword
    from keywords, movie_keywords
    where movie_keywords.movie_id IN (`+strings.Join(qp, ",")+`)
    and movie_keywords.keyword_id = keywords.id`, args...,
  )
  if err != nil {
    log.Println("Error querying postgres!")
    log.Fatal(err)
  }

  var keywords []string
  for rows.Next() {
    var keyword string
    if err := rows.Scan(&keyword); err != nil {
      log.Fatal(err)
    }
    keywords = appendIfMissing(keywords, keyword)
  }
  return keywords
}

func appendIfMissing(slice []string, s string) []string {
  for _, ele := range slice {
    if ele == s {
      return slice
    }
  }
  return append(slice, s)
}
