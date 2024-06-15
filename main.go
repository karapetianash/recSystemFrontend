package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	r := gin.Default()
	db, err := sql.Open("sqlite3", "../data/recommendations.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r.GET("/recommendations/:userId/:n", func(c *gin.Context) {
		getRecommendationsHandler(c, db)
	})

	r.Run(":8080")
}
