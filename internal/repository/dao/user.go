package dao

import (
	"control/internal/errno"
	"control/internal/repository/model"
	"database/sql"
	"github.com/denisenkom/go-mssqldb"
	"github.com/pkg/errors"
)

func (d *database) IsAdmin(userId string) (bool, error) {
	sqlSentence := "SELECT is_admin FROM library_user WHERE user_id = @user_id"
	row := d.db.QueryRow(sqlSentence, sql.Named("user_id", userId))
	
	var isAdmin bool
	if err := row.Scan(&isAdmin); err != nil {
		return false, errors.WithStack(err)
	}
	
	return isAdmin, nil
}

func (d *database) UserHistory(borrow *model.Borrow) ([]*model.BorrowBook, error) {
	sqlSentence := "SELECT borrow_id,borrow_time,should_return_time,return_time," +
		"book.book_id,book_name,book_author,book_publish_time,book_used,user_id " +
		"FROM borrow JOIN book ON borrow.book_id = book.book_id " +
		"AND borrow.user_id = @user_id"
	
	rows, err := d.db.Query(sqlSentence, sql.Named("user_id", borrow.UserId))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	
	ret := make([]*model.BorrowBook, 0)
	for rows.Next() {
		var one model.BorrowBook
		if err := rows.Scan(&one.BorrowId, &one.BorrowTime, &one.ShouldReturnTime, &one.ReturnTime,
			&one.Book.BookId, &one.BookName, &one.BookAuthor, &one.BookPublishedTime, &one.BookUsed,
			&one.UserId);
			err != nil {
			return ret, errors.WithStack(err)
		}
		ret = append(ret, &one)
	}
	return ret, nil
}

func (d *database) UserRegister(user *model.User) error {
	sqlSentence := "INSERT INTO library_user(user_id,user_password,user_name,is_admin) " +
		"VALUES(@user_id,@user_password,@user_name,@is_admin)"
	if _, err := d.db.Exec(sqlSentence, sql.Named("user_id", user.UserId),
		sql.Named("user_password", user.UserPassword), sql.Named("user_name", user.UserName),
		sql.Named("is_admin", user.IsAdmin)); err != nil {
		if sqlErr, ok := err.(mssql.Error); ok && sqlErr.Number == 2601 || sqlErr.Number == 2627 {
			return errno.ErrUserExist
		}
		return errors.WithStack(err)
	}
	return nil
}

func (d *database) UserLogin(user *model.User) error {
	sqlSentence := "SELECT user_id,user_password FROM library_user WHERE user_id = @user_id"
	row := d.db.QueryRow(sqlSentence, sql.Named("user_id", user.UserId))
	
	var one model.User
	if err := row.Scan(&one.UserId, &one.UserPassword); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errno.ErrUserNotExist
		}
		return errors.WithStack(err)
	}
	if user.UserPassword != one.UserPassword {
		return errno.ErrWrongPassword
	}
	
	return nil
}

func (d *database) ReturnBook(borrow *model.Borrow) error {
	sqlSentence := "SELECT book_used FROM book JOIN borrow ON " +
		"book.book_id=borrow.book_id AND borrow_id=@borrow_id"
	row := d.db.QueryRow(sqlSentence, sql.Named("borrow_id", borrow.BorrowId))
	
	var used bool
	if err := row.Scan(&used); err != nil {
		return errors.WithStack(err)
	}
	
	if !used {
		return errno.ErrSpareBook
	}
	
	sqlSentence = "DELETE FROM borrow WHERE borrow_id = @borrow_id"
	
	if _, err := d.db.Exec(sqlSentence, sql.Named("borrow_id",
		borrow.BorrowId)); err != nil {
		return errors.WithStack(err)
	}
	
	return nil
}
