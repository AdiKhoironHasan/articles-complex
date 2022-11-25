package database

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type SqlDatabase struct {
	Conn *sqlx.DB
}

type NoSqlDatabase struct {
	Conn *redis.Client
}

func PostgreSqllInitialize(host, username, password, database, port string) (SqlDatabase, error) {
	db := SqlDatabase{}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	fmt.Println("dsn postgresql:", dsn)
	conn, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return db, err
	}

	db.Conn = conn
	err = db.Conn.Ping()

	if err != nil {
		return db, err
	}

	return db, nil
}
func RedislInitialize(host, username, password, port string) (NoSqlDatabase, error) {
	db := NoSqlDatabase{}

	address := fmt.Sprintf("%s:%s", host, port)
	fmt.Println("dsn redis:", address)
	db.Conn = redis.NewClient(&redis.Options{
		Addr: address,
		// Username: username,
		Password: password,
		DB:       0,
	})

	err := db.Conn.Ping(context.Background()).Err()
	if err != nil {
		return db, err
	}

	return db, nil
}
