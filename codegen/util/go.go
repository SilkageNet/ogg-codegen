package util

var goNumTypes = []string{
	"float32",
	"float64",
	"int",
	"int8",
	"int16",
	"int32",
	"int64",
	"uint",
	"uint8",
	"uint16",
	"uint32",
	"uint64"}

func IsGoNumType(t string) bool {
	return StrInArray(t, goNumTypes)
}

func IsGoStrType(t string) bool {
	return t == "string"
}
