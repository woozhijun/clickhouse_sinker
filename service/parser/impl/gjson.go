package impl

import (
	"github.com/tidwall/gjson"
	"github.com/woozhijun/clickhouse_sinker/service/parser"
	"github.com/woozhijun/clickhouse_sinker/util"
	"time"
)

func NewParser(typ string) parser.Parser {
	switch typ {
	case "json", "gjson":
		return &GjsonParser{}
	default:
		return &GjsonParser{}
	}
}

type GjsonParser struct {
}

func (c *GjsonParser) Parse(bs []byte) parser.IParser {
	return &GjsonMetric{string(bs)}
}

type GjsonMetric struct {
	raw string
}

func (c *GjsonMetric) Get(key string) interface{} {
	return gjson.Get(c.raw, key).Value()
}

func (c *GjsonMetric) GetString(key string) string {
	return gjson.Get(c.raw, key).String()
}

func (c *GjsonMetric) GetFloat(key string) float64 {
	return gjson.Get(c.raw, key).Float()
}

func (c *GjsonMetric) GetInt(key string) int64 {
	return gjson.Get(c.raw, key).Int()
}

func (c *GjsonMetric) GetDate(key string, layout string) string {
	if key == "_CURRENT_" {
		return time.Now().UTC().Format(layout)
	}
	date := util.StringParseToDate(gjson.Get(c.raw, key).String(), layout)
	return date.In(time.FixedZone("UTC", -8*60*60)).Format(util.LayoutDatetime)
}