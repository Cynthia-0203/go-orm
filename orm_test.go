package orm

import (
	"errors"
	"fmt"
	"reflect"

	"testing"

	"github.com/Cynthia/go-orm/session"
	_ "github.com/go-sql-driver/mysql"
)

func OpenDB(t *testing.T) *Engine {
	t.Helper()
	engine, err := NewEngine("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v","root","123","localhost","3306","orm"))
	if err != nil {
		t.Fatal("failed to connect", err)
	}
	return engine
}

func TestNewEngine(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
}

type User struct {
	Name string `orm:"PRIMARY KEY"`
	Age  int
}



func transactionRollback(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	_ = s.Model(&User{}).DropTable()
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		_ = s.Model(&User{}).CreateTable()
		_, _ = s.Insert(&User{"Tom", 18})
		return nil, errors.New("Error")
	})
	if err == nil {
		t.Fatal("failed to rollback")
	}
}
// func transactionCommit(t *testing.T) {
// 	engine := OpenDB(t)
// 	defer engine.Close()
// 	s := engine.NewSession()
// 	_ = s.Model(&User{}).DropTable()
// 	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
// 		_ = s.Model(&User{}).CreateTable()
// 		_, err = s.Insert(&User{"Tom", 18})
// 		return
// 	})
// 	u := &User{}
// 	_ = s.First(u)
// 	if err != nil || u.Name != "Tom" {
// 		t.Fatal("failed to commit")
// 	}
// }

func TestEngine_Transaction(t *testing.T) {
	t.Run("rollback", func(t *testing.T) {
		transactionRollback(t)
	})
	// t.Run("commit", func(t *testing.T) {
	// 	transactionCommit(t)
	// })
}

func TestEngine_Migrate(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name varchar PRIMARY KEY, XXX int);").Exec()
	_, _ = s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	engine.Migrate(&User{})

	rows, _ := s.Raw("SELECT * FROM User").QueryRows()
	columns, _ := rows.Columns()
	if !reflect.DeepEqual(columns, []string{"Name", "Age"}) {
		t.Fatal("Failed to migrate table User, got columns", columns)
	}
}