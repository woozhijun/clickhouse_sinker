package data_type

import (
	"fmt"
)

type IntColumn struct {
	name string
}

func NewIntColumn(bits int, isUint bool) *IntColumn {
	name := fmt.Sprintf("Int%d", bits)
	if isUint {
		name = "U" + name
	}
	return &IntColumn{name: name}
}

func (c *IntColumn) Name() string {
	return c.name
}

func (c *IntColumn) DefaultValue() interface{} {
	return int64(0)
}

// only judge int and float64
func (c *IntColumn) GetValue(val interface{}) interface{} {
	switch val.(type) {
	case int:
		return int64(val.(int))
	case float64:
		return int64(val.(float64))
	}
	return int64(0)
}