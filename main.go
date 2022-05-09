package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity uint   `json:"quantity"`
}

var books []book = []book{
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 5},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 5},
}

// GET /books
func getBooks(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, books)
}

// POST /books
func createBook(context *gin.Context) {
	var newBook book

	if err := context.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)
	context.IndentedJSON(http.StatusCreated, newBook)
}

// GET /books/:id
func bookById(context *gin.Context) {
	id := context.Param("id")
	book, err := getBookById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}
	context.IndentedJSON(http.StatusOK, book)
}

// PUT /checkout?id=:id
func checkoutBook(context *gin.Context) {
	id, ok := context.GetQuery("id")

	// error: problem retrieving the id
	if !ok {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	// success: try to retrieve the book with the id
	book, err := getBookById(id)

	// error: book not found
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	// error: book not available
	if book.Quantity <= 0 {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available"})
	}

	// success: book checked out
	book.Quantity -= 1
	context.IndentedJSON(http.StatusOK, book)
}

// PUT /return?id=:id
func returnBook(context *gin.Context) {
	id, ok := context.GetQuery("id")

	if !ok {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing or bad id"})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	book.Quantity += 1
	context.IndentedJSON(http.StatusOK, book)
}

// utitlity
func getBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}
	return nil, errors.New("book not found")
}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", bookById)
	router.PUT("/return", returnBook)
	router.PUT("/checkout", checkoutBook)
	router.POST("/books", createBook)

	router.Run("localhost:8080")
}
