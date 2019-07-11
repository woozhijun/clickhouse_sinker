package model

type ItemSchema struct {
	Name 	string
	Type 	string
	Alias   string
	Layout string
}

type MetricItem struct {
	Name string
	Type string
}

func (itemSchema *ItemSchema) ChooseAliasName() string {
	if itemSchema.Alias != "" {
		return itemSchema.Alias
	}
	return itemSchema.Name
}