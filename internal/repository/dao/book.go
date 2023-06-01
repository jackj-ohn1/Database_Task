package dao

import (
	"control/internal/errno"
	"control/internal/repository/model"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"time"
)

func (d *database) GetBookId(bookName string) (int, error) {
	sqlSentence := "SELECT book_id FROM book WHERE book_name = @book_name"
	
	var bookId int
	row := d.db.QueryRow(sqlSentence, sql.Named("book_name", bookName))
	if err := row.Scan(&bookId); err != nil {
		return -1, errors.WithStack(err)
	}
	
	return bookId, nil
}

func (d *database) GetBooks(page, limit int) ([]*model.Book, error) {
	sqlSentence := "SELECT book_id,book_name,book_author,book_publish_time,book_used,book_out FROM book ORDER BY book_used DESC OFFSET " +
		"@offset ROWS FETCH NEXT @limit ROWS ONLY"
	rows, err := d.db.Query(sqlSentence, sql.Named("offset", (page-1)*limit),
		sql.Named("limit", limit))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	
	var data = make([]*model.Book, 0)
	for rows.Next() {
		var book model.Book
		err := rows.Scan(&book.BookId, &book.BookName, &book.BookAuthor,
			&book.BookPublishedTime, &book.BookUsed, &book.BookOut)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		
		data = append(data, &book)
	}
	
	return data, nil
}

func (d *database) AddBook(book *model.Book) error {
	sqlSentence := "INSERT INTO book(book_name,book_author,book_publish_time) " +
		"VALUES(@name,@author,@time)"
	
	if _, err := d.db.Exec(sqlSentence, sql.Named("name", book.BookName),
		sql.Named("author", book.BookAuthor), sql.Named("time",
			book.BookPublishedTime)); err != nil {
		return err
	}
	
	return nil
}

func (d *database) DeleteBook(book *model.Book) error {
	sqlSentence := "DELETE FROM book WHERE book_id = @book_id AND book_used = 0"
	if _, err := d.db.Exec(sqlSentence, sql.Named("book_id", book.BookId)); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (d *database) LendBook(borrow *model.Borrow, borrowDay int) error {
	sqlSentence := "SELECT book_used FROM book WHERE book_id=@book_id"
	row := d.db.QueryRow(sqlSentence, sql.Named("book_id", borrow.BookId))
	
	var used bool
	var max, borrowed int
	if err := row.Scan(&used); err != nil {
		return errors.WithStack(err)
	}
	
	if used {
		return errno.ErrUsedBook
	}
	
	sqlSentence = "SELECT borrowed_book,borrow_max FROM library_user WHERE user_id = @user_id"
	row = d.db.QueryRow(sqlSentence, sql.Named("user_id", borrow.UserId))
	if err := row.Scan(&borrowed, &max); err != nil {
		return errors.WithStack(err)
	}
	if borrowed >= max {
		return errno.ErrNoLeftResource
	}
	
	duration, err := time.ParseDuration(fmt.Sprintf(
		"%dh", borrowDay*24))
	if err != nil {
		return errors.WithStack(err)
	}
	borrow.ShouldReturnTime = time.Now().Add(duration)
	
	sqlSentence = "INSERT INTO borrow(user_id, book_id, should_return_time)" +
		"VALUES(@user_id, @book_id, @should_return_time)"
	
	if _, err := d.db.Exec(sqlSentence, sql.Named("user_id", borrow.UserId),
		sql.Named("book_id", borrow.BookId),
		sql.Named("should_return_time", borrow.ShouldReturnTime)); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
