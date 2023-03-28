package codegen

type EnumValue struct {
	Name  string
	Value string
}

type EnumValues []EnumValue

func (e EnumValues) Len() int {
	return len(e)
}

func (e EnumValues) Less(i, j int) bool {
	return e[i].Value < e[j].Value
}

func (e EnumValues) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}
