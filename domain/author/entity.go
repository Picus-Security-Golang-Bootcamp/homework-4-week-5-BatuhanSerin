package author

import (
	"fmt"

	"github.com/BatuhanSerin/postgresql/domain/book"
	"gorm.io/gorm"
)

type Author struct {
	gorm.Model
	AuthorName string
	AuthorID   string
	Books      []book.Book `gorm:"foreignkey:AuthorID;references:AuthorID"`
}
type authorSlice []Author

// ToString returns author information
func (a Author) ToString() string {
	return fmt.Sprintf("\nAuthor id: %s\nAuthor name: %s", a.AuthorID, a.AuthorName)
}
