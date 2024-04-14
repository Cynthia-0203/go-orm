package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/Cynthia/go-orm"
)

func main(){

	engine,_:=orm.NewEngine("mysql",fmt.Sprintf("%v:%v@tcp(%v:%v)/%v","root","123","localhost","3306","orm"))
	defer engine.Close()

	s:=engine.NewSession()
	_,_=s.Raw("drop table if exists users;").Exec()
	_,_=s.Raw("create table users(Name text);").Exec()
	result,_:=s.Raw("insert into users(`Name`) values (?),(?)","Tom","ja").Exec()
	count,_:=result.RowsAffected()
	fmt.Printf("Exce success,%d affected\n",count)
}

