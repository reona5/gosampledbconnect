package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dbconn := "root:password@tcp(localhost:3306)/go_sample?parseTime=true"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := Open(dbconn)
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	Ping(db, ctx)
	GetFirstUser(db, ctx)
}

func Open(dbconn string) *sql.DB {
	db, err := sql.Open("mysql", dbconn)
	if err != nil {
		log.Panicf("sql.Open() failed. error: %v", err.Error())
	}
	return db
}

func Ping(db *sql.DB, ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Panicf("db.PingContext() failed:%v\n", err.Error())
	}
}

func GetFirstUser(db *sql.DB, ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	type User struct {
		UserID string `db:"id"`
		UserName string `db:"name"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}
	var (
		id, name             string
		createdAt, updatedAt time.Time
	)
	var firstUser *User
	query := "SELECT * FROM users WHERE id = 1;"
	row := db.QueryRowContext(ctx, query)

	row.Scan(&id, &name, &createdAt, &updatedAt)
	firstUser = &User{
		UserID: id,
		UserName: name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	fmt.Println(firstUser)
}

