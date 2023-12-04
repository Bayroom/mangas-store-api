package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// manga (struct) represents data about a record manga-book
type manga struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Rating float64 `json:"rating"` // rating/10
}

// mangas (slice) represents a record sample of mangas
var mangaStore = []manga{
	{ID: "1", Title: "Hunter x Hunter", Author: "Yoshihiro Togashi", Rating: 9.8},
	{ID: "2", Title: "Jujutsu Kaisen", Author: "Gege Akutami", Rating: 8.9},
	{ID: "3", Title: "Death Note", Author: "Obata Takeshi", Rating: 8.6},
	{ID: "4", Title: "Dr. Stone", Author: "Riichiro Inagaki", Rating: 9.5},
}

// welcome represents a func to welcome visitors
func welcome(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Welcome to our humble mangas store !\n")
}

// getMangas represents a func to retrieve mangaStore in json format
func getMangas(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, mangaStore)
}

// getMangaById represents a func to retrieve a single element from mangaStore in json format
func getMangaById(ctx *gin.Context) {
	// Retrieve the ID from the URL parameters
	id := ctx.Param("id")

	// Loop over mangaStore to return manga if present
	for _, manga := range mangaStore {
		if manga.ID == id {
			ctx.IndentedJSON(http.StatusOK, manga)
			return
		}

	}

	// this line reached no manga found
	ctx.IndentedJSON(http.StatusNotFound, gin.H{"Error": "No manga with that ID is stored in our mangaStore.."})
}

// addManga represents a func to add a new element to mangaStore
func addManga(ctx *gin.Context) {
	var newManga manga

	// bind recieved manga to newManga var
	if err := ctx.BindJSON(&newManga); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Bad request format"})
		return
	}

	// add newManga to slice
	mangaStore = append(mangaStore, newManga)
	ctx.IndentedJSON(http.StatusCreated, mangaStore)
}

// updateManga represents a func to update an element already in mangaStore
func updateManga(ctx *gin.Context) {
	// Retrieve the ID from the URL parameters
	id := ctx.Param("id")

	// Find the index of the manga with the specified ID
	index := -1
	for idx, manga := range mangaStore {
		if manga.ID == id {
			index = idx
			break
		}
	}

	// If the manga with the specified ID is not found, return an error
	if index == -1 {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "Manga not found"})
		return
	}

	// Bind the updated manga data from the request body
	var updatedManga manga
	if err := ctx.ShouldBindJSON(&updatedManga); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the existing manga in the store
	mangaStore[index] = updatedManga

	ctx.IndentedJSON(http.StatusOK, updatedManga)

}

// deleteManga represents a func to delete an element already in mangaStore
func deleteManga(ctx *gin.Context) {
	// Retrieve the ID from the URL parameters
	id := ctx.Param("id")

	// Find the index of the manga with the specified ID
	index := -1
	for idx, manga := range mangaStore {
		if manga.ID == id {
			index = idx
			break
		}
	}

	// If the manga with the specified ID is not found, return an error
	if index == -1 {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "Manga not found"})
		return
	}

	// Remove the manga from the store
	mangaStore = append(mangaStore[:index], mangaStore[index+1:]...)
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Manga deleted successfully", "mangaStore": mangaStore})
}

func main() {
	// Default
	router := gin.Default()

	// Possible routes
	router.GET("/", welcome)
	router.GET("/mangas", getMangas)
	router.GET("/mangas/:id", getMangaById)
	router.POST("/mangas", addManga)
	router.PUT("/mangas/:id", updateManga)
	router.DELETE("/mangas/:id", deleteManga)

	// Run the APIs
	router.Run(":8080")
}
