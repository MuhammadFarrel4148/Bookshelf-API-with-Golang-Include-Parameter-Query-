package handler

import (
	"bookshelfapi/books"
	"bookshelfapi/database"
	"bookshelfapi/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateBookInput struct {
	Name      string `json:"name" binding:"required"`
	Year      int    `json:"year" binding:"required"`
	Author    string `json:"author" binding:"required"`
	Summary   string `json:"summary" binding:"required"`
	Publisher string `json:"publisher" binding:"required"`
	PageCount int    `json:"pageCount" binding:"required"`
	ReadPage  int    `json:"readPage" binding:"required"`
	Reading   int    `json:"reading" binding:"required"`
}


func SavedBooks(c *gin.Context) {
	var input CreateBookInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Gagal menambahkan buku. Mohon isi nama buku",
		})
		return
	}

	if input.ReadPage > input.PageCount {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Gagal menambahkan buku. readPage tidak boleh lebih besar dari pageCount",
		})
		return
	}

	book := books.BookInput{
		ID:         utils.GenerateID(),
		Name:       input.Name,
		Year:       input.Year,
		Author:     input.Author,
		Summary:    input.Summary,
		Publisher:  input.Publisher,
		PageCount:  input.PageCount,
		ReadPage:   input.ReadPage,
		Finished:   input.ReadPage == input.PageCount,
		Reading:    input.Reading == 1,
		InsertedAt: time.Now(),
		UpdatedAt:  time.Now(),
	}

	database.DB.Create(&book)

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Buku berhasil ditambahkan",
		"data": gin.H{
			"bookId": book.ID,
		},
	})
}

func GetAllBooks(c *gin.Context) {
	var books []books.BookInput
	readingQuery := c.Query("reading")
	finishedQuery := c.Query("finished")

	if readingQuery != "" {
		var isReading bool

		if readingQuery == "1" {
			isReading = true
		} else {
			isReading = false
		}

		result := database.DB.Where("reading = ?", isReading).Find(&books).Error
		if result != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"status": "fail",
				"message": "Gagal mendapatkan buku",
			})
			return
		}
	}

	if finishedQuery != "" {
		result := database.DB.Where("finished = ?", finishedQuery).Find(&books).Error
		if result != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "fail",
				"message": "Gagal mendapatkan buku",
			})
			return
		}
	}
	
	err := database.DB.Find(&books).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"message": "Gagal menampilkan buku",
		})
		return
	}

	var response []gin.H
	for _, book := range books {
		response = append(response, gin.H{
			"id": book.ID,
			"name": book.Name,
			"publisher": book.Publisher,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"books": response,
		},
	})
}

func GetBookById(c *gin.Context) {
	bookId := c.Param("bookId")

	var book books.BookInput
	err := database.DB.First(&book, "id = ?", bookId).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "fail",
			"message": "Buku tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"book": gin.H{
				"id": book.ID,
				"name": book.Name,
				"year": book.Year,
				"author": book.Author,
				"summary": book.Summary,
				"publisher": book.Publisher,
				"pageCount": book.PageCount,
				"readPage": book.ReadPage,
				"finished": book.Finished,
				"reading": book.Reading,
				"insertedAt": book.InsertedAt,
				"updatedAt": book.UpdatedAt,
			},
		},
	})
}

func UpdateBook(c *gin.Context) {
	bookId := c.Param("bookId")

	var book books.BookInput
	result := database.DB.First(&book, "id = ?", bookId).Error

	if result != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "fail",
			"message": "Gagal memperbarui buku. Id tidak ditemukan",
		})
		return
	}

	var books books.BookInput
	err := c.ShouldBindJSON(&books)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"message": "Gagal memperbarui buku. Mohon isi nama buku",
		})
		return
	}

	if books.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"message": "Gagal memperbarui buku. Mohon isi nama buku",
		})
		return
	}

	if books.ReadPage > books.PageCount {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"message": "Gagal memperbarui buku. readPage tidak boleh lebih besar dari pageCount",
		})
		return
	}

	book.Name = books.Name
	book.Year = books.Year
	book.Author = books.Author
	book.Summary = books.Summary
	book.Publisher = books.Publisher
	book.PageCount = books.PageCount
	book.ReadPage = books.ReadPage
	book.Reading = books.Reading

	database.DB.Save(&book)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "Buku berhasil diperbarui",
	})
}

func DeleteBookById(c *gin.Context) {
	bookId := c.Param("bookId")

	var book books.BookInput

	err := database.DB.First(&book, "id = ?", bookId).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "fail",
			"message": "Buku gagal dihapus. Id tidak ditemukan",
		})
		return
	}
	
	result := database.DB.Delete(&book).Error

	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"message": "Gagal menghapus buku",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "Buku berhasil dihapus",
	})	
}
