package dao

import (
	"control/internal/errno"
	"control/internal/repository/model"
	"database/sql"
	"github.com/denisenkom/go-mssqldb"
	"github.com/pkg/errors"
)

func (d *database) UserHistory() {

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

func (d *database) ReturnBook() {

}
