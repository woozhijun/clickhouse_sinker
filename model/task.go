package model


type Task struct {
	Name string

	Kafka         string
	Topic         string
	ConsumerGroup string

	// Earliest set to true to consume the message from oldest position
	Earliest bool
	Parser   string

	Clickhouse string
	TableName  string

	// AutoSchema will auto fetch the schema from clickhouse
	AutoSchema     bool
	ExcludeColumns []string

	ItemSchemas []ItemSchema  `json:"itemSchemas"`
	Metrics []MetricItem `json:"metrics"`

	FlushInterval int `json:"flushInterval,omitempty"`
	BufferSize    int `json:"bufferSize,omitempty"`
}
