package creator

import (
	_ "github.com/kshvakov/clickhouse"
	"github.com/prometheus/common/log"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/woozhijun/clickhouse_sinker/util"
	"testing"
	"time"
)

var (
	c *Config
)

func TestInitConfig(t *testing.T) {
	Convey("init config", t, func() {
		c = InitConfig("../conf")

	})
}

func TestConfig_GenTasks(t *testing.T) {

	Convey("check tasks input", t, func() {
		c = InitConfig("../conf")
		for _, taskService := range c.GenTasks() {
			taskService.Init()
			taskService.Run()
		}
	})
}

func TestConfig_GenInput(t *testing.T) {

	Convey("check tasks input", t, func() {
		c = InitConfig("../conf")
		//out := make(chan []byte, 30000)
		for _, task := range c.Tasks {
			kafka := c.GenInput(task)
			kafka.Init()
		}
	})
}

func TestDate(t *testing.T) {

	Convey("check output datetime", t, func() {

		d := util.StringParseToDate("2019-07-04T19:00:00+08:00", "2006-01-02T15:04:05+08:00")
		println(d.In(time.FixedZone("UTC", -8*60*60)).Format(util.LayoutDatetime))
		println(d.UTC().Format(util.LayoutDatetime))
		log.Info(time.Now().UTC().Format(util.LayoutDatetime))
		log.Info(time.Now().In(time.FixedZone("CST", 8*60*60)).Format(util.LayoutDatetime))
	})
}