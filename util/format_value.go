package util

import (
	"github.com/woozhijun/clickhouse_sinker/model"
	"github.com/woozhijun/clickhouse_sinker/service/parser"
)

//这里对metric的value类型，只有三种情况， （float64，string，map[string]interface{})
func GetValueByType(metric parser.IParser, cwt *model.ItemSchema) interface{} {
	swType := switchType(cwt.Type)
	name := cwt.ChooseAliasName()
	switch swType {
	case "int":
		return metric.GetInt(name)
	case "float":
		return metric.GetFloat(name)
	case "string":
		return metric.GetString(name)
	case "date":
		return metric.GetDate(name, cwt.Layout)
	//never happen
	default:
		return ""
	}
}

func switchType(typ string) string {
	switch typ {
	case "UInt8", "UInt16", "UInt32", "UInt64", "Int8", "Int16", "Int32", "Int64":
		return "int"
	case "String", "FixString":
		return "string"
	case "Float32", "Float64":
		return "float"
	case "Date", "DateTime":
		return "date"
	default:
		panic("unsupport type " + typ)
	}
}