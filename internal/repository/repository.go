package repository

import (
	"control/internal/repository/dao"
	"control/internal/repository/model"
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
)

type Database interface {
	UserRegister(user *model.User) error
	UserLogin(user *model.User) error
}

func NewDatabase(db *sql.DB) Database {
	return dao.NewDatabase(db)
}

func Init(username, password string) *sql.DB {
	db, err := sql.Open("sqlserver",
		fmt.Sprintf("sqlserver://%s:%s@localhost:%d?database=book_control",
			username, password, 1433))
	if err != nil {
		log.Fatal("数据库连接失败!")
	}
	return db
}
