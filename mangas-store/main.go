package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// manga (struct) represents data about a record manga-book
type manga struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Rating float64 `json:"rating"` // rating/10
}

// createTableQuery is the query that adds mangas table to the postgres mangas database
var createTableQuery = `
CREATE TABLE IF NOT EXISTS mangaStore (
	id		SERIAL PRIMARY KEY,
	title	TEXT,
	author	TEXT,
	rating	FLOAT
);
`

// initialDataQuery is the query to insert initial data into the mangaStore table
var initialDataQuery = `
INSERT INTO mangaStore (title, author, rating) VALUES
    ('Hunter x Hunter', 'Yoshihiro Togashi', 9.8),
    ('Jujutsu Kaisen', 'Gege Akutami', 8.9),
    ('Death Note', 'Obata Takeshi', 8.6),
    ('Dr. Stone', 'Riichiro Inagaki', 9.5);
`

// Check if the table is empty
func isTableEmpty(db *sql.DB) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM mangaStore").Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// initDB initializes the DB and creates the table(s)
func initDB() (*sql.DB, error) {
	// PostgreSQL connection string for Dockerized database
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Create mangas table if it does not exist
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	// Check if the table is empty before running initialDataQuery
	isEmpty, err := isTableEmpty(db)
	if err != nil {
		return nil, err
	}

	// Insert initial data into mangaStore table if table is empty
	if isEmpty {
		_, err = db.Exec(initialDataQuery)
		if err != nil {
			return nil, err
		}

	}

	return db, nil
}

// welcome represents a func to welcome visitors
func welcome(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Welcome to our humble mangas store !\n")
}

// getMangas represents a func to retrieve mangaStore in json format
func getMangas(ctx *gin.Context, db *sql.DB) {
	rows, err := db.Query("SELECT * FROM mangaStore")
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error:": "Internal server Error"})
		return
	}
	defer rows.Close()

	var mangas []manga
	for rows.Next() {
		var m manga
		if err := rows.Scan(&m.ID, &m.Title, &m.Author, &m.Rating); err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Internal server Error"})
			return
		}
		mangas = append(mangas, m)
	}

	// Return mangas
	ctx.IndentedJSON(http.StatusOK, mangas)
}

// getMangaById represents a func to retrieve a single element from mangaStore in json format
func getMangaById(ctx *gin.Context, db *sql.DB) {
	// Retrieve the ID from the URL parameters
	id := ctx.Param("id")

	var m manga
	err := db.QueryRow("SELECT * FROM mangaStore WHERE id=$1", id).Scan(&m.ID, &m.Title, &m.Author, &m.Rating)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.IndentedJSON(http.StatusNotFound, gin.H{"Error": "Manga not found"})
			return
		}
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Internal server Error"})
		return
	}

	// Return manga
	ctx.IndentedJSON(http.StatusOK, m)
}

// addManga represents a func to add a new element to mangaStore
func addManga(ctx *gin.Context, db *sql.DB) {
	var newManga manga

	// bind recieved manga to newManga var
	if err := ctx.BindJSON(&newManga); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Bad request format"})
		return
	}

	// add newManga to database
	_, err := db.Exec("INSERT INTO mangaStore (title, author, rating) VALUES ($1, $2, $3)",
		newManga.Title, newManga.Author, newManga.Rating)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Internal server Error"})
		return
	}
	ctx.IndentedJSON(http.StatusCreated, gin.H{"message": "Manga added successfully"})
}

// updateManga represents a func to update an element already in mangaStore
func updateManga(ctx *gin.Context, db *sql.DB) {
	// Retrieve the ID from the URL parameters
	id := ctx.Param("id")

	// Find the manga with the specified ID
	var m manga
	err := db.QueryRow("SELECT * FROM mangaStore WHERE id=$1", id).Scan(&m.ID, &m.Title, &m.Author, &m.Rating)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.IndentedJSON(http.StatusNotFound, gin.H{"Error": "Manga not found"})
			return
		}
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Internal server Error"})
		return
	}

	// Bind the updated manga data from the request body
	var updatedManga manga
	if err := ctx.ShouldBindJSON(&updatedManga); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the existing manga in mangaStore
	_, err = db.Exec("UPDATE mangaStore SET title=$1, author=$2, rating=$3 WHERE id=$4",
		updatedManga.Title, updatedManga.Author, updatedManga.Rating, id)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Internal server Error"})
		return
	}

	// Retuen updated manga
	ctx.IndentedJSON(http.StatusOK, updatedManga)
}

// deleteManga represents a func to delete an element already in mangaStore
func deleteManga(ctx *gin.Context, db *sql.DB) {
	// Retrieve the ID from the URL parameters
	id := ctx.Param("id")

	// Find the manga with the specified ID
	var m manga
	err := db.QueryRow("SELECT * FROM mangaStore WHERE id=$1", id).Scan(&m.ID, &m.Title, &m.Author, &m.Rating)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.IndentedJSON(http.StatusNotFound, gin.H{"Error": "Manga not found"})
			return
		}
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Internal server Error"})
		return
	}

	// Remove the manga from the mangaStore
	_, err = db.Exec("DELETE FROM mangaStore WHERE id=$1", id)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Internal server Error"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Manga deleted successfully"})
}

func main() {
	// Initialize the database
	db, err := initDB()
	if err != nil {
		panic(err)
	}
	// Close DB when done
	defer db.Close()

	// Default
	router := gin.Default()

	// Possible routes
	router.GET("/", welcome)
	router.GET("/mangas", func(ctx *gin.Context) {
		getMangas(ctx, db)
	})
	router.GET("/mangas/:id", func(ctx *gin.Context) {
		getMangaById(ctx, db)
	})
	router.POST("/mangas", func(ctx *gin.Context) {
		addManga(ctx, db)
	})
	router.PUT("/mangas/:id", func(ctx *gin.Context) {
		updateManga(ctx, db)
	})
	router.DELETE("/mangas/:id", func(ctx *gin.Context) {
		deleteManga(ctx, db)
	})

	// Run the APIs
	router.Run(":8080")
}
