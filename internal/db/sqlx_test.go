package db_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
)

func TestIn(t *testing.T) {
	q, args, err := sqlx.In("(?) (?)", []any{[]int{1, 2}, []int{3, 4}}...)
	fmt.Printf("q: %#v\nargs:%#v\nerr:%v\n", q, args, err)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

const ISO8601 = "2006-01-02T15:04:05.000Z"

type ISO8601Time string

func encode(t time.Time) *ISO8601Time {
	i := ISO8601Time(t.UTC().Format(ISO8601))
	return &i
}

func TestTime(t *testing.T) {
	db := sqlx.MustConnect("sqlite3", ":memory:")
	defer db.Close()
	db.MustExec("CREATE TABLE foo (time TEXT)")
	x, err := db.NamedQuery("INSERT INTO foo (time) VALUES (:time) RETURNING time", []struct {
		Time *ISO8601Time `db:"time"`
	}{
		{encode(time.Now())},
		{encode(time.Date(2020, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))},
	})
	must(err)
	for x.Next() {
		var t ISO8601Time
		must(x.Scan(&t))
		fmt.Printf("t: %#v\n", t)
	}
	rows, err := db.Queryx("SELECT * FROM foo")
	must(err)
	// rows.Nextre()
	for rows.Next() {
		var t ISO8601Time
		must(rows.Scan(&t))
		fmt.Printf("t: %#v\n", t)
	}
}

func TestFormat(t *testing.T) { 
	fmt.Printf("%#v\n", time.Now().UTC().Format(time.RFC3339Nano))
}