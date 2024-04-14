package dialect

import (
	"fmt"
	"reflect"
	"time"
)

type mysql struct{}

var _ Dialect=(*mysql)(nil)

func init(){
	RegisterDialect("mysql",&mysql{})
}

func (s *mysql) DataTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Bool:
		return "bool"
	case  reflect.Int8, 
		reflect.Uint, reflect.Uint8:
		return "smallint"
	case reflect.Int16, reflect.Int32,reflect.Uint16, reflect.Uint32:
		return "int"
	case reflect.Int64, reflect.Uint64,reflect.Int,reflect.Uintptr:
		return "bigint"
	case reflect.Float32:
		return "float"
	case reflect.Float64:
		return "double"
	case reflect.String:
		return "varchar"
	// case reflect.Array, reflect.Slice:
	// 	return "blob"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))
}


func(s *mysql)TableExistSQL(tableName string)(string,[]interface{}){
	args:=[]interface{}{tableName}
	return "SELECT table_name FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?", args
}