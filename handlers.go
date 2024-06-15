package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Recommendation struct {
	MovieId         int     `json:"movie_id"`
	PredictedRating float64 `json:"predicted_rating"`
}

func getRecommendationsHandler(c *gin.Context, db *sql.DB) {
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userId"})
		return
	}

	n, err := strconv.Atoi(c.Param("n"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid number of recommendations"})
		return
	}

	if !userExists(db, userId) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	recommendations, err := getRecommendations(db, userId, n)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, recommendations)
}

func userExists(db *sql.DB, userId int) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM ratings WHERE userId=? LIMIT 1)`
	err := db.QueryRow(query, userId).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false
	}
	return exists
}

func getRecommendations(db *sql.DB, userId int, n int) ([]Recommendation, error) {
	query := `
        SELECT movieId, predicted_rating
        FROM recommendations
        WHERE userId=?
        ORDER BY predicted_rating DESC
        LIMIT ?
    `
	rows, err := db.Query(query, userId, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	recommendations := make([]Recommendation, 0)
	for rows.Next() {
		var rec Recommendation
		if err := rows.Scan(&rec.MovieId, &rec.PredictedRating); err != nil {
			return nil, err
		}
		recommendations = append(recommendations, rec)
	}
	return recommendations, nil
}
