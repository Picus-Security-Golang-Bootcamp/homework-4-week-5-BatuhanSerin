package author

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"gorm.io/gorm"
)

//AuthorRepository is a struct for AuthorRepository
type AuthorRepository struct {
	db *gorm.DB
}

//NewAuthorRepository returns Author Repository
func NewAuthorRepository(db *gorm.DB) *AuthorRepository {
	return &AuthorRepository{db: db}
}

//GetAllAuthorsWithBookInformation returns all authors with book information
func (a *AuthorRepository) GetAllAuthorsWithBookInformation() (authorSlice, error) {
	var authors authorSlice
	result := a.db.Preload("Books").Find(&authors)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, author := range authors {
		fmt.Println(author.ToString())
		fmt.Println("Books: ")
		if len(author.Books) > 0 {
			for _, book := range author.Books {
				fmt.Println(book.ToString())
				fmt.Println("=============================")
			}
		}
	}
	return authors, nil
}

//GetAuthorWithName returns author by its name
func (a *AuthorRepository) GetAuthorWithName(name string) error {
	var authors *Author
	result := a.db.Where(Author{AuthorName: name}).Preload("Books").Find(&authors)
	if result.Error != nil {
		return result.Error
	}

	fmt.Println(authors.ToString())
	fmt.Println("Books: ")
	if len(authors.Books) > 0 {
		for _, book := range authors.Books {
			fmt.Println(book.ToString())
			fmt.Println("=============================")
		}
	}

	return nil
}

//**********************************______________________********************
//Migrations Auto Migrates for authors
func (a *AuthorRepository) Migrations() {
	a.db.AutoMigrate(&Author{})
}

//InsertData inserts data from csv file to database with ReadCsvAuthor function
func (a *AuthorRepository) InsertData() {

	err := a.ReadCsvAuthor()
	if err != nil {
		log.Fatal(err)
	}

}

//ReadCsvAuthor reads datas from csv file
func (a *AuthorRepository) ReadCsvAuthor() error {
	f, err := os.Open("author.csv")
	if err != nil {
		return err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return err
	}

	var authors = authorSlice{}
	for _, line := range records[1:] {
		authors = append(authors, Author{
			AuthorID:   line[0],
			AuthorName: line[1],
		})
	}
	for _, author := range authors {
		a.db.Where(Author{AuthorName: author.AuthorName}).
			Attrs(Author{AuthorID: author.AuthorID, AuthorName: author.AuthorName}).
			FirstOrCreate(&author)
	}
	return nil
}
