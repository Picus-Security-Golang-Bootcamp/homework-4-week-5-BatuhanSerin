package book

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

//BookRepository is a struct for BookRepository
type BookRepository struct {
	db *gorm.DB
}

//NewBookRepository returns Book Repository
func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

// FindAll returns all informations of books
func (b *BookRepository) FinAll() bookSlice {
	var books bookSlice
	b.db.Find(&books)
	fmt.Println("Books: ")
	if len(books) > 0 {
		for _, book := range books {
			fmt.Println(book.ToString())
			fmt.Println("=============================")
		}

	}
	return books
}

//FindBookById returns book by its ID
func (b *BookRepository) FindBookById(id int) bookSlice {
	var books bookSlice
	strID := strconv.Itoa(id)
	//b.db.Where("id = ?", strID).Order("id desc , name").Find(&books)
	b.db.Where(&Book{ID: strID}).Order("id desc , name").Find(&books)
	fmt.Println("Books: ")
	if len(books) > 0 {
		for _, book := range books {
			fmt.Println(book.ToString())
			fmt.Println("=============================")
		}
	}
	return books
}

//FindByAuthorOrBookId returns book by its author id or book id
func (b *BookRepository) FindByAuthorOrBookId(id int) bookSlice {
	var books bookSlice
	strID := strconv.Itoa(id)

	b.db.Where("id = ?", strID).Or("author_id = ?", strID).Find(&books)
	fmt.Println("Books: ")
	if len(books) > 0 {
		for _, book := range books {
			fmt.Println(book.ToString())
			fmt.Println("=============================")
		}
	}
	return books

}

//FindByName returns book by its name
func (b *BookRepository) FindByName(name string) bookSlice {
	var books bookSlice
	Name := strings.Title(strings.ToLower(name))
	b.db.Where("name LIKE ? ", "%"+Name+"%").Find(&books)
	fmt.Println("Books: ")
	if len(books) > 0 {
		for _, book := range books {
			fmt.Println(book.ToString())
			fmt.Println("=============================")
		}
	}
	return books
}

//FindByNameWithRawSql returns book by its name with raw sql
func (b *BookRepository) FindByNameWithRawSql(name string) bookSlice {
	var books bookSlice
	b.db.Raw("SELECT * FROM books WHERE name LIKE ? ", "%"+name+"%").Scan(&books)

	fmt.Println("Books: ")
	if len(books) > 0 {
		for _, book := range books {
			fmt.Println(book.ToString())
			fmt.Println("=============================")
		}
	}

	return books
}

//GetByID returns book by its ID
func (b *BookRepository) GetByID(id int) (*Book, error) {

	var book Book
	strID := strconv.Itoa(id)

	result := b.db.First(&book, strID)
	fmt.Println("Book: ")

	fmt.Println(book.ToString())
	fmt.Println("=============================")

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	return &book, nil
}

//Create creates book in database
func (b *BookRepository) Create(book *Book) error {
	result := b.db.Create(book)
	if result.Error != nil {
		return result.Error
	}
	return nil

}

//Update updates book in database
func (b *BookRepository) Update(book *Book) error {
	result := b.db.Save(book)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

//Delete deletes book from database
func (b *BookRepository) Delete(book *Book) error {
	result := b.db.Delete(book)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

//DeleteById deletes book by its ID from database without checking the book is deleted or not
func (b *BookRepository) DeleteById(id int) error {
	strID := strconv.Itoa(id)
	book := Book{ID: strID}
	result := b.db.Delete(&book)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

//BeforeDelete deletes book from database after checking the book is deleted or not
func (b *BookRepository) BeforeDelete(id int) (err error) {
	strID := strconv.Itoa(id)
	var book Book
	result := b.db.First(&book, strID)
	if result.Error != nil {
		return result.Error
	}

	if book.DeletedAt.Time.IsZero() {
		fmt.Println("Deleted Book: ")

		fmt.Println(book.ToString())
		fmt.Println("=============================")

		result := b.db.Delete(&book)
		if result.Error != nil {
			return errors.New("This book has already been deleted")
		}
		return nil
	}
	return nil
}

//********************************************_____________________________*************************************
//Migrations Auto Migrates for books
func (b *BookRepository) Migrations() {
	b.db.AutoMigrate(&Book{})
}

//InsertData inserts data from csv file to database with ReadCsvBook function
func (b *BookRepository) InsertData() {
	err := b.ReadCsvBook()
	if err != nil {
		log.Fatal(err)
	}
}

//ReadCsvBook reads datas from csv file
func (b *BookRepository) ReadCsvBook() error {
	f, err := os.Open("book.csv")
	if err != nil {
		return err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return err
	}

	var books = bookSlice{}
	for _, line := range records[1:] {
		books = append(books, Book{
			ID:        line[0],
			Name:      line[1],
			Page:      line[2],
			Stock:     line[3],
			Cost:      line[4],
			StockCode: line[5],
			ISBN:      line[6],
			AuthorID:  line[7],
		})
	}
	for _, book := range books {
		b.db.Where(Book{Name: book.Name}).
			Attrs(Book{ID: book.ID, Name: book.Name}).
			FirstOrCreate(&book)
	}
	return nil
}
