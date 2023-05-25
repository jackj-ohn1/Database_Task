package model

import "time"

type Book struct {
	BookId            int     `json:"book_id"`
	BookName          string  `json:"book_name"`
	BookAuthor        string  `json:"book_author"`
	BookPublishedTime string  `json:"book_published_time"`
	BookUsed          int     `json:"book_used"`
	BookPrice         float32 `json:"book_price"`
	BookContent       string  `json:"book_content"`
}

type User struct {
	UserId       string `json:"user_id"`
	UserPassword string `json:"user_password"`
	UserName     string `json:"user_name"`
	IsAdmin      int    `json:"is_admin"`
	BorrowMax    int    `json:"borrow_max"`
	BorrowedBook int    `json:"borrowed_book"`
}

type Borrow struct {
	BorrowId         int       `json:"borrow_id"`
	BookId           int       `json:"book_id"`
	UserId           int       `json:"user_id"`
	BorrowTime       time.Time `json:"borrow_time"`
	ShouldReturnTime time.Time `json:"should_return_time"`
	ReturnTime       time.Time `json:"return_time"`
}
