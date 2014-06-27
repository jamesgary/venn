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
  psqlHost := os.Getenv("SNOWFLAKE_PSQL_HOST")
  assetsDir := os.Getenv("SNOWFLAKE_ASSETS_DIR")

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

  server.Get("/guess/:keywords", func(params martini.Params, w http.ResponseWriter) string {
    w.Header().Set("Content-Type", "application/json")

    keywordIds := getKeywordIdsFromStrings(strings.Split(params["keywords"], " "), db)
    var movieIdsForKeyword [][]int
    for _, keywordId := range keywordIds {
      movieIdsForKeyword = append(movieIdsForKeyword, getMovieIdsWithKeywordId(keywordId, db))
    }
    intersectedMovieIds := getIntersection(movieIdsForKeyword)

    if len(intersectedMovieIds) == 1 {
      title := getMovieFromId(intersectedMovieIds[0], db)
      return fmt.Sprintf(`{"movie": %s}`, strconv.Quote(title))
    } else {
      possibleKeywords := getKeywordsForMovies(intersectedMovieIds, db)
      sort.Strings(possibleKeywords)
      for i, keyword := range possibleKeywords {
        possibleKeywords[i] = strconv.Quote(keyword)
      }
      return fmt.Sprintf(`{"movieCount": %d, "possibleKeywords": [%s]}`, len(intersectedMovieIds), strings.Join(possibleKeywords, ","))
    }
  })

  //server.Run()
  log.Fatal(http.ListenAndServe(":3001", server))
}

func getMovieFromId(id int, db *sql.DB) string {
  rows, err := db.Query("select title from movies where id = $1", id)
  if err != nil {
    log.Println("Error querying postgres!")
    log.Fatal(err)
  }
  var movie string
  for rows.Next() {
    if err := rows.Scan(&movie); err != nil {
      log.Fatal(err)
    }
  }
  return movie
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

func getKeywordIdsFromStrings(keywords []string, db *sql.DB) []int {
  args := make([]interface{}, len(keywords))
  qp := make([]string, len(keywords))
  for i, v := range keywords {
    args[i] = interface{}(v)
    qp[i] = fmt.Sprintf("$%d", i+1)
  }
  rows, err := db.Query(
    `select id from keywords where keyword IN (`+strings.Join(qp, ",")+`)`,
    args...,
  )
  if err != nil {
    log.Println("Error querying postgres!")
    log.Fatal(err)
  }

  var keywordIds []int
  for rows.Next() {
    var keywordId int
    if err := rows.Scan(&keywordId); err != nil {
      log.Fatal(err)
    }
    keywordIds = append(keywordIds, keywordId)
  }
  return keywordIds
}

func getIntersection(idSets [][]int) []int {
  var intersection []int
  for _, id := range idSets[0] {
    if isInAllArrays(id, idSets) {
      intersection = append(intersection, id)
    }
  }
  return intersection
}

func isIn(i int, set []int) bool {
  for _, k := range set {
    if i == k {
      return true
    }
  }
  return false
}

func isInAllArrays(i int, sets [][]int) bool {
  for _, k := range sets {
    if !isIn(i, k) {
      return false
    }
  }
  return true
}
