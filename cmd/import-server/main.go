package main

import (
    "log"
    "net/http"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "github.com/rtsf-ratings/parser"
    "github.com/rtsf-ratings/ratings/db"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

type Player struct {
    Id int

    Firstname string
    Lastname  string
}

type ImportData struct {
    Players    map[int]*db.Player
    PlayersMap map[int]Player
    Result     parser.Result
}

func main() {
    router := gin.Default()
    router.Use(cors.Default())

    dsn := mysql.Open("ratings:ratings@tcp(127.0.0.1:3336)/ratings_new?charset=utf8&parseTime=true")
    ratings, err := db.RatingsOpen(dsn, &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    router.GET("/players", func(c *gin.Context) {
        c.IndentedJSON(http.StatusOK, ratings.Players)
    })

    router.POST("/players", func(c *gin.Context) {
        var data ImportData

        if err := c.BindJSON(&data); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "unknown file type"})
            return
        }

        c.IndentedJSON(http.StatusOK, gin.H{"status": "success"})
    })

    router.Run("127.0.0.1:8083")
}
