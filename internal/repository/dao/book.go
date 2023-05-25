package dao

import (
	"control/internal/repository/model"
	"fmt"
	"github.com/pkg/errors"
	"time"
)

func (d *database) GetBooks(page, limit int) ([]*model.Book, error) {
	sqlSentence := fmt.Sprintf("SELECT * FROM book OFFSET %d ROWS FETCH NEXT %d ROWS ONLY",
		(page-1)*limit, limit)
	rows, err := d.db.Query(sqlSentence)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	
	var data = make([]*model.Book, 0)
	for rows.Next() {
		var book model.Book
		err := rows.Scan(&book.BookId, &book.BookName, &book.BookAuthor,
			&book.BookPublishedTime, &book.BookUsed, &book.BookPrice, &book.BookContent)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		
		data = append(data, &book)
	}
	
	return data, nil
}

func (d *database) AddBook() {

}

func (d *database) DeleteBook(bookId int) error {
	sqlSentence := fmt.Sprintf("DELETE FROM book WHERE book_id = %d", bookId)
	if _, err := d.db.Exec(sqlSentence); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (d *database) LendBook(userId, bookId, borrowTime int) error {
	duration, err := time.ParseDuration(fmt.Sprintf(
		"%dh", borrowTime*24))
	if err != nil {
		return errors.WithStack(err)
	}
	shouldReturnTime := time.Now().Add(duration)
	
	sqlSentence := "INSERT INTO borrow(user_id, book_id, should_return_time)" +
		"VALUES(?, ?, ?)"
	
	if _, err := d.db.Exec(sqlSentence, userId, bookId, shouldReturnTime); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
