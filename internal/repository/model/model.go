package model

import (
	"database/sql"
	"time"
)

type Book struct {
	BookId            int          `json:"book_id,omitempty"`
	BookName          string       `json:"book_name,omitempty"`
	BookAuthor        string       `json:"book_author,omitempty"`
	BookPublishedTime sql.NullTime `json:"book_published_time,omitempty"`
	BookUsed          bool         `json:"book_used,omitempty"`
	//BookPrice         float32 `json:"book_price,omitempty"`
	//BookContent       string  `json:"book_content,omitempty"`
}

type User struct {
	UserId       string `json:"user_id,omitempty"`
	UserPassword string `json:"user_password,omitempty"`
	UserName     string `json:"user_name,omitempty"`
	IsAdmin      bool   `json:"is_admin,omitempty"`
	BorrowMax    int    `json:"borrow_max,omitempty"`
	BorrowedBook int    `json:"borrowed_book,omitempty"`
}

type Borrow struct {
	BorrowId         int          `json:"borrow_id,omitempty"`
	BookId           int          `json:"book_id,omitempty"`
	UserId           string       `json:"user_id,omitempty"`
	BorrowTime       time.Time    `json:"borrow_time,omitempty"`
	ShouldReturnTime time.Time    `json:"should_return_time,omitempty"`
	ReturnTime       sql.NullTime `json:"return_time,omitempty"`
}

type BorrowBook struct {
	Borrow
	Book
}
