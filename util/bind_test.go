package util

import (
	"fmt"
	"github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"github.com/wswz/go_commons/utils"
	"testing"
)

func TestBind(t *testing.T) {
	config := map[string]interface{}{
		"a": "A",
		"b": 1,
		"c": 3.0,
		"x": 1,
		"z": "zz",
	}

	type Foo struct {
		X int
		Z string
	}
	type Entity struct {
		A string
		B int
		C float32
		Foo
	}
	var entity = &Entity{}
	IngestConfig(config, entity)

	assert.Equal(t, &Entity{A: "A", B: 1, C: 3.0, Foo: Foo{1, "zz"}}, entity, "must be equal")
}

func Test(t *testing.T)  {
	convey.Convey("test ip4 By name ", t, func() {

		var hosts []string
		ips, err := utils.GetIp4Byname("redash")
		if err != nil {
			println(err)
		}
		for _, ip := range ips {
			hosts = append(hosts, fmt.Sprintf("%s:%d", ip, 9000))
		}

		for i, host := range hosts {
			println("host:" + string(i) +  " " + host)
		}
	})
}