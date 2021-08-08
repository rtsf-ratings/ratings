package main

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
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

    router.GET("/tournament/:id", func(c *gin.Context) {
        var data ImportData

        id, err := uuid.Parse(c.Param("name"))
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
            return
        }

        content, err := ioutil.ReadFile("uploads/" + id.String())
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "unknown uuid"})
            return
        }

        err = json.Unmarshal([]byte(content), &data)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
            return
        }

        c.IndentedJSON(http.StatusOK, &data)
    })

    router.POST("/tournament", func(c *gin.Context) {
        var data ImportData

        if err := c.BindJSON(&data); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "unknown file type"})
            return
        }

        json, err := json.MarshalIndent(data, "", "  ")
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
            return
        }

        id := uuid.New()

        err = ioutil.WriteFile("uploads/"+id.String(), json, 0644)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
            return
        }

        c.IndentedJSON(http.StatusOK, gin.H{"status": "success", "id": id.String()})
    })

    router.POST("/tournament/:id", func(c *gin.Context) {
        var data ImportData

        if err := c.BindJSON(&data); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "unknown file type"})
            return
        }

        json, err := json.MarshalIndent(data, "", "  ")
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
            return
        }

        id, err := uuid.Parse(c.Param("name"))
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
            return
        }

        err = ioutil.WriteFile("uploads/"+id.String(), json, 0644)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
            return
        }

        c.IndentedJSON(http.StatusOK, gin.H{"status": "success", "id": id.String()})
    })

    router.Run("127.0.0.1:8083")
}
