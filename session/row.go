package session

import (
	"database/sql"
	"strings"

	"github.com/Cynthia/go-orm/clause"
	"github.com/Cynthia/go-orm/dialect"
	"github.com/Cynthia/go-orm/logger"
	"github.com/Cynthia/go-orm/schema"
)

type Session struct{
	db *sql.DB
	sql strings.Builder
	sqlVars []interface{}

	dialect dialect.Dialect
	refTable *schema.Schema

	clause clause.Clause

	tx *sql.Tx
}
func New(db *sql.DB,dialect dialect.Dialect) *Session {
	return &Session{
		db:db,
		dialect: dialect,
	}
}

func (s *Session) Clear(){
	s.sql.Reset()
	s.sqlVars=nil
	s.clause=clause.Clause{}
}

type CommonDB interface{
	Query(query string,args ...interface{})(*sql.Rows,error)
	QueryRow(query string,args ...interface{})*sql.Row
	Exec(query string,args ...interface{})(sql.Result,error)
}


var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)

func (s *Session)DB()CommonDB{
	if s.tx!=nil{
		return s.tx
	}
	return s.db
}



func(s *Session)Raw(sql string,values ...interface{})*Session{
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars=append(s.sqlVars, values...)
	return s
}

func(s *Session)Exec()(result sql.Result,err error){
	defer s.Clear()
	logger.Info(s.sql.String(),s.sqlVars)
	if result,err=s.DB().Exec(s.sql.String(),s.sqlVars...);err!=nil{
		logger.Error(err)
	}
	return
}

func(s *Session)QueryRow()*sql.Row{
	defer s.Clear()

	logger.Info(s.sql.String(),s.sqlVars)
	return s.DB().QueryRow(s.sql.String(),s.sqlVars...)
}

func(s *Session)QueryRows()(rows *sql.Rows,err error){
	defer s.Clear()
	logger.Info(s.sql.String(),s.sqlVars)

	if rows,err=s.DB().Query(s.sql.String(),s.sqlVars...);err!=nil{
		logger.Error(err)
	}

	return
}
