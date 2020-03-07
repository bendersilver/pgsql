package pgsql

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/bendersilver/blog"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgxpool"
)

var dbURL string

// Init -
func Init() {
	dbURL = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASS"),
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		os.Getenv("PG_DB"),
	)
}

// PGPool -
func PGPool() (*pgxpool.Pool, error) {
	return pgxpool.Connect(context.Background(), dbURL)
}

// Select  -
type Select struct {
	sql  string
	args []interface{}
}

// Item -
func (s *Select) Item(fields ...interface{}) error {
	pool, err := PGPool()
	if err != nil {
		return err
	}
	defer pool.Close()
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	return conn.QueryRow(context.Background(), s.sql, s.args...).Scan(fields...)
}

// Get -
func Get(sql string, args ...interface{}) *Select {
	s := new(Select)
	s.sql = sql
	s.args = args
	return s
}

// Set -
func Set(sql string, arg ...interface{}) error {
	pool, err := PGPool()
	if err != nil {
		return err
	}
	defer pool.Close()
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	_, err = conn.Exec(context.Background(), sql, arg...)
	return err
}

// Batch -
func Batch(sql string, args [][]interface{}) error {
	pool, err := PGPool()
	if err != nil {
		return err
	}
	defer pool.Close()
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	batch := &pgx.Batch{}
	for i, arg := range args {
		batch.Queue(sql, arg...)
		args[i] = nil
	}
	br := conn.SendBatch(context.Background(), batch)
	if err := br.Close(); err != nil {
		return err
	}
	return nil
}

var notifyFunc []func(string)

// AddNotify -
func AddNotify(f func(string)) {
	notifyFunc = append(notifyFunc, f)
}

func pgnotify(name string) error {
	pool, err := PGPool()
	if err != nil {
		return err
	}
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	_, err = conn.Exec(context.Background(), fmt.Sprintf("listen %s", name))
	for {
		notify, err := conn.Conn().WaitForNotification(context.Background())
		if err != nil {
			return err
		}
		for _, n := range notifyFunc {
			n(notify.Payload)
		}
	}
	return nil

}

// RunNotify -
func RunNotify(name string) {
	go func() {
		for {
			err := pgnotify(name)
			blog.Error("reconnecr pgNotify", err)
			time.Sleep(time.Second * 2)
		}
	}()
}
