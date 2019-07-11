package column

import (
	"github.com/woozhijun/clickhouse_sinker/column"
	"github.com/woozhijun/clickhouse_sinker/column/data_type"
)

var (
	columns = map[string]column.IColumn{}
)

type creator func() column.IColumn

func registor(name string, creator creator) {
	columns[name] = creator()
}

// GetColumnByName get the IColumn by the name of type
func GetColumnByName(name string) column.IColumn {
	return columns[name]
}

func init() {
	registor("UInt8", func() column.IColumn {
		return data_type.NewIntColumn(8, false)
	})
	registor("UInt16", func() column.IColumn {
		return data_type.NewIntColumn(16, false)
	})
	registor("UInt32", func() column.IColumn {
		return data_type.NewIntColumn(32, false)
	})
	registor("UInt64", func() column.IColumn {
		return data_type.NewIntColumn(64, false)
	})

	registor("Int8", func() column.IColumn {
		return data_type.NewIntColumn(8, false)
	})
	registor("Int16", func() column.IColumn {
		return data_type.NewIntColumn(16, false)
	})
	registor("Int32", func() column.IColumn {
		return data_type.NewIntColumn(32, false)
	})
	registor("Int64", func() column.IColumn {
		return data_type.NewIntColumn(64, false)
	})

	registor("Float32", func() column.IColumn {
		return data_type.NewFloatColumn(32)
	})
	registor("Float64", func() column.IColumn {
		return data_type.NewFloatColumn(64)
	})

	registor("String", func() column.IColumn {
		return data_type.NewStringColumn()
	})

	registor("FixedString", func() column.IColumn {
		return data_type.NewStringColumn()
	})

	registor("Date", func() column.IColumn {
		return data_type.NewDateColumn( false)
	})

	registor("DateTime", func() column.IColumn {
		return data_type.NewDateColumn( true)
	})
}