package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func setupRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()
	r.GET("/recommendations/:userId/:n", func(c *gin.Context) {
		getRecommendationsHandler(c, db)
	})
	return r
}

func setupTestDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	createTables(db)
	insertTestData(db)
	return db, nil
}

func createTables(db *sql.DB) {
	createTablesSQL := `
    CREATE TABLE IF NOT EXISTS ratings (
        userId INTEGER,
        movieId INTEGER,
        rating REAL,
        PRIMARY KEY (userId, movieId)
    );
    CREATE TABLE IF NOT EXISTS recommendations (
        userId INTEGER,
        movieId INTEGER,
        predicted_rating REAL,
        PRIMARY KEY (userId, movieId)
    );
    `
	_, err := db.Exec(createTablesSQL)
	if err != nil {
		panic(err)
	}
}

func insertTestData(db *sql.DB) {
	insertRatingsSQL := `
    INSERT INTO ratings (userId, movieId, rating) VALUES (1, 1, 4.0);
    INSERT INTO ratings (userId, movieId, rating) VALUES (1, 2, 5.0);
    INSERT INTO ratings (userId, movieId, rating) VALUES (2, 1, 3.0);
    INSERT INTO ratings (userId, movieId, rating) VALUES (2, 2, 2.0);
    INSERT INTO ratings (userId, movieId, rating) VALUES (2, 3, 5.0);
    `
	_, err := db.Exec(insertRatingsSQL)
	if err != nil {
		panic(err)
	}

	insertRecommendationsSQL := `
    INSERT INTO recommendations (userId, movieId, predicted_rating) VALUES (1, 3, 4.5);
    INSERT INTO recommendations (userId, movieId, predicted_rating) VALUES (1, 2, 4.2);
    INSERT INTO recommendations (userId, movieId, predicted_rating) VALUES (2, 1, 3.5);
    `
	_, err = db.Exec(insertRecommendationsSQL)
	if err != nil {
		panic(err)
	}
}

func TestGetRecommendationsHandler(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}
	defer db.Close()

	router := setupRouter(db)

	t.Run("Valid userId and n", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/recommendations/1/2", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		expectedBody := `[{"movie_id":3,"predicted_rating":4.5},{"movie_id":2,"predicted_rating":4.2}]`
		if strings.TrimSpace(w.Body.String()) != expectedBody {
			t.Fatalf("Expected body %s, got %s", expectedBody, w.Body.String())
		}
	})

	t.Run("Invalid userId", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/recommendations/abc/2", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}

		expectedBody := `{"error":"Invalid userId"}`
		if strings.TrimSpace(w.Body.String()) != expectedBody {
			t.Fatalf("Expected body %s, got %s", expectedBody, w.Body.String())
		}
	})

	t.Run("Non-existent userId", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/recommendations/999/2", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
		}

		expectedBody := `{"error":"User not found"}`
		if strings.TrimSpace(w.Body.String()) != expectedBody {
			t.Fatalf("Expected body %s, got %s", expectedBody, w.Body.String())
		}
	})

	t.Run("Invalid number of recommendations", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/recommendations/1/abc", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}

		expectedBody := `{"error":"Invalid number of recommendations"}`
		if strings.TrimSpace(w.Body.String()) != expectedBody {
			t.Fatalf("Expected body %s, got %s", expectedBody, w.Body.String())
		}
	})
}
