package session

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/Cynthia/go-orm/dialect"

	_ "github.com/go-sql-driver/mysql"
)

var (
	TestDB      *sql.DB
	TestDial, _ = dialect.GetDialect("mysql")
)

func TestMain(m *testing.M) {
	TestDB, _ = sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v","root","123","localhost","3306","orm"))
	code := m.Run()
	_ = TestDB.Close()
	os.Exit(code)
}

func NewSession() *Session {
	return New(TestDB, TestDial)
}
func TestSession_Exec(t *testing.T) {
	s := NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	result, _ := s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	if count, err := result.RowsAffected(); err != nil || count != 2 {
		t.Fatal("expect 2, but got", count)
	}
}

func TestSession_QueryRows(t *testing.T) {
	s := NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	row := s.Raw("SELECT count(*) FROM User").QueryRow()
	var count int
	if err := row.Scan(&count); err != nil || count != 0 {
		t.Fatal("failed to query db", err)
	}
}