package data_type

import "github.com/woozhijun/clickhouse_sinker/util"

type DateColumn struct {
	name string
	layout string
}

func NewDateColumn(isTime bool) *DateColumn {
	if isTime {
		return &DateColumn{"DateTime", util.LayoutDatetime}
	}
	return &DateColumn{"Date", util.LayoutDate}
}

func (c *DateColumn) Name() string {
	return c.name
}

func (c *DateColumn) DefaultValue() interface{} {
	return ""
}

func (c *DateColumn) GetValue(val interface{}) interface{} {
	switch val.(type) {
	case string:
		return val.(string)
	}
	return ""
}