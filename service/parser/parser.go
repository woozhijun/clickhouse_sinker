package parser

type IParser interface {
	Get(key string) interface{}
	GetString(key string) string
	GetFloat(key string) float64
	GetInt(key string) int64
	GetDate(key string, layout string) string
}

type Parser interface {
	Parse(bs []byte) IParser
}