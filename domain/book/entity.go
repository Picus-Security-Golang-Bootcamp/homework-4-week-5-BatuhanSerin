package book

import (
	"fmt"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	ID        string
	Name      string
	Page      string
	Stock     string
	Cost      string
	StockCode string
	ISBN      string
	AuthorID  string
}

type bookSlice []Book

// ToString returns book information
func (book Book) ToString() string {
	return fmt.Sprintf("id: %s\nName: %s\nPage: %s\nStock: %s\nCost: %s\nStockCode: %s\nISBN: %s",
		book.ID, book.Name, book.Page, book.Stock, book.Cost, book.StockCode, book.ISBN)
}
