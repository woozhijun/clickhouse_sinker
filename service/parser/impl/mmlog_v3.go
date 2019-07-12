package impl

import (
	"github.com/tidwall/gjson"
	"github.com/woozhijun/clickhouse_sinker/service/parser"
	"github.com/woozhijun/clickhouse_sinker/util"
	"regexp"
	"time"
)

var (
	regex = "(\\S*) (\\S*) (\\S*) (\\S*) (\\S*) (\\S*) (\\S*) (\\S*) (\\S*) (\\{.*\\}$)"
	re = regexp.MustCompile(regex)
)

type MMlogV3Parser struct {
}

func (c *MMlogV3Parser) Parse(bs []byte) parser.IParser {
	baseInfo := make(map[string]string)
	baseInfo["severity"] = "$1"
	baseInfo["datetime"] = "$2"
	baseInfo["host"] = "$3"
	baseInfo["service_name"] = "$4"
	baseInfo["process_id"] = "$5"
	baseInfo["thread_id"] = "$6"
	baseInfo["log_version"] = "$7"
	baseInfo["query_id"] = "$8"
	baseInfo["event_name"] = "$9"
	baseInfo["raw"] = "$10"
	template,_ := json.Marshal(baseInfo)
	if re.Match(bs) {
		m := re.FindSubmatchIndex(bs)
		return &MMlogV3Metric{string(re.Expand([]byte{}, template, bs, m)), baseInfo["raw"]}
	} else {
		return &MMlogV3Metric{"", string(bs)}
	}
}

type MMlogV3Metric struct {
	baseInfo		string
	raw 			string
}

func (c *MMlogV3Metric) chooseGjsonResult(key string) gjson.Result {
	if c.baseInfo != "" && gjson.Get(c.baseInfo, key).String() != "" {
		return gjson.Get(c.baseInfo, key)
	}
	return gjson.Get(c.raw, key)
}


func (c *MMlogV3Metric) Get(key string) interface{} {
	return c.chooseGjsonResult(key).Value()
}

func (c *MMlogV3Metric) GetString(key string) string {
	return c.chooseGjsonResult(key).String()
}

func (c *MMlogV3Metric) GetFloat(key string) float64 {
	return c.chooseGjsonResult(key).Float()
}

func (c *MMlogV3Metric) GetInt(key string) int64 {
	return c.chooseGjsonResult(key).Int()
}

func (c *MMlogV3Metric) GetDate(key string, layout string) string {
	if key == "_CURRENT_" {
		return time.Now().UTC().Format(layout)
	}
	date := util.StringParseToDate(gjson.Get(c.raw, key).String(), layout)
	return date.In(time.FixedZone("UTC", -8*60*60)).Format(util.LayoutDatetime)
}