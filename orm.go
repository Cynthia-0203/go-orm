package orm

import (
	"database/sql"
	"fmt"

	"github.com/Cynthia/go-orm/dialect"
	"github.com/Cynthia/go-orm/logger"
	"github.com/Cynthia/go-orm/session"
)

type Engine struct{
	db *sql.DB
	dialect dialect.Dialect
}

type TxFunc func (*session.Session) (interface{},error)

func (engine *Engine)Transaction(f TxFunc)(result interface{},err error){
	s:=engine.NewSession()
	if s==nil{
		return nil,err
	}
	if err:=s.Begin();err!=nil{
		return nil,err
	}

	defer func(){
		if p:=recover();p!=nil{
			_=s.Rollback()
			panic(p)
		}else if err!=nil{
			_=s.Rollback()

		}else{
			defer func ()  {
				if err != nil {
					_ = s.Rollback()
				}
			}()
			err=s.Commit()
		}
	}()

	return f(s)
}
func NewEngine(driver,source string)(e *Engine,err error){
	db,err:=sql.Open(driver,source)
	if err!=nil{
		logger.Error(err)
		return
	}

	if err=db.Ping();err!=nil{
		logger.Error(err)
		return
	}

	dial,ok:=dialect.GetDialect(driver)

	if !ok{
		logger.Errorf("dialect %s not found",driver)
		return
	}


	e=&Engine{db: db,dialect:dial}

	logger.Info("Connected database successfully...")
	return
}

func(engine *Engine)Close(){
	if err:=engine.db.Close();err!=nil{
		logger.Error("Failed to close database...")
	}
	logger.Info("Close database successfully...")
}

func(engine *Engine)NewSession()*session.Session{
	return session.New(engine.db,engine.dialect)
}

func difference(a []string, b []string) (diff []string) {
	mapB := make(map[string]bool)
	for _, v := range b {
		mapB[v] = true
	}
	for _, v := range a {
		if _, ok := mapB[v]; !ok {
			diff = append(diff, v)
		}
	}
	return

}

func(engin *Engine)Migrate(value interface{})error{
	_,err:=engin.Transaction(func(s *session.Session)(result interface{}, err error){
		if !s.Model(value).HasTable(){
			logger.Infof("table %s doesn't exist",s.RefTable().Name)
			return nil,s.CreateTable()
		}

		table:=s.RefTable()
		rows,_:=s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", table.Name)).QueryRows()
		columns,_:=rows.Columns()
		addCols:=difference(table.FieldNames,columns)
		delCols:=difference(columns,table.FieldNames)
		logger.Infof("added cols %v, deleted cols %v", addCols, delCols)

		for _, col := range addCols {
			f := table.GetField(col)
			sqlStr := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", table.Name, f.Name, f.Type)
			if _, err = s.Raw(sqlStr).Exec(); err != nil {
				return
			}
		}
		if len(delCols) == 0 {
			return
		}

		for _, col := range delCols {
			f := table.GetField(col)
			sqlStr := fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s;", table.Name, f.Name)
			if _, err = s.Raw(sqlStr).Exec(); err != nil {
				return
			}
		}

		return
	})
	return err
}