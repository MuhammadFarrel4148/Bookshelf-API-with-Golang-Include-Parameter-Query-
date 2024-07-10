package main

import (
	"bookshelfapi/database"
	"bookshelfapi/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	database.ConnectDatabase()

	router.POST("/books", handler.SavedBooks)
	router.GET("/books", handler.GetAllBooks)
	router.GET("/books/:bookId", handler.GetBookById)
	router.PUT("/books/:bookId", handler.UpdateBook)
	router.DELETE("/books/:bookId", handler.DeleteBookById)

	router.Run()
}
