package main

import (
	"log"

	postgres "github.com/BatuhanSerin/postgresql/common/db"
	"github.com/BatuhanSerin/postgresql/domain/author"
	"github.com/BatuhanSerin/postgresql/domain/book"

	//bookStruct "github.com/BatuhanSerin/postgresql/domain/book"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := postgres.NewPsqlDB()
	if err != nil {
		log.Fatal("Postgres cannot init ", err)
	}

	log.Println("Postgres connected")

	bookRepo := book.NewBookRepository(db)
	bookRepo.Migrations()
	bookRepo.InsertData()

	//fmt.Println(bookRepo.FinAll())
	//fmt.Println(bookRepo.FindBookById(2))
	//fmt.Println(bookRepo.FindByAuthorOrBookId(5))
	//fmt.Println(bookRepo.FindByName("It"))
	//fmt.Println(bookRepo.FindByNameWithRawSql("It"))
	// fmt.Println(bookRepo.GetByID(2))

	// NewBook := bookStruct.Book{ID: "4", Name: "The Dice Man", Page: "305", Stock: "14", Cost: "25", StockCode: "7", ISBN: "A125-128-DCD", AuthorID: "7"}
	// bookRepo.Create(&NewBook)
	// fmt.Println(bookRepo.FinAll())
	// bookRepo.Delete(&NewBook)
	// fmt.Println(bookRepo.FinAll())
	//bookRepo.BeforeDelete(2)

	// Author************

	authorRepo := author.NewAuthorRepository(db)
	authorRepo.Migrations()
	authorRepo.InsertData()

	//authorRepo.GetAuthorWithName("Jack London")
	//authorRepo.GetAllAuthorsWithBookInformation()

}