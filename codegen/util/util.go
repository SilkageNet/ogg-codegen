package util

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"strings"
	"text/template"
	"unicode"
)

// ToCamelCase This function will convert query-arg style strings to CamelCase. We will
// use `., -, +, :, ;, _, ~, ' ', (, ), {, }, [, ]` as valid delimiters for words.
// So, "word.word-word+word:word;word_word~word word(word)word{word}[word]"
// would be converted to WordWordWordWordWordWordWordWordWordWordWordWordWord
func ToCamelCase(str string) string {
	var separators = "-#@!$&=.+:;_~ (){}[]"
	var s = strings.Trim(str, " ")
	var n = ""
	var capNext = true
	for _, v := range s {
		if unicode.IsUpper(v) {
			n += string(v)
		}
		if unicode.IsDigit(v) {
			n += string(v)
		}
		if unicode.IsLower(v) {
			if capNext {
				n += strings.ToUpper(string(v))
			} else {
				n += string(v)
			}
		}
		capNext = strings.ContainsRune(separators, v)
	}
	return n
}

// IsInternalRef 判断是否为内部引用类型
func IsInternalRef(ref string) bool {
	return ref != "" && strings.HasPrefix(ref, "#")
}

// ConvRef2GoType 将ref转成golang类型名
func ConvRef2GoType(ref string) (string, error) {
	if !IsInternalRef(ref) {
		return "", fmt.Errorf("ref.unsupported:%s", ref)
	}
	var pathParts = strings.Split(ref, "/")
	if depth := len(pathParts); depth != 4 {
		return "", fmt.Errorf("ref.unsupported:%s", ref)
	}
	return ToCamelCase(pathParts[3]), nil
}

// Str2GoComment renders a possible multi-line string as a valid Go-Comment.
// Each line is prefixed as a comment.
func Str2GoComment(in string, name string) string {
	if len(in) == 0 || len(strings.TrimSpace(in)) == 0 { // ignore empty comment
		return ""
	}

	// Normalize newlines from Windows/Mac to Linux
	in = strings.Replace(in, "\r\n", "\n", -1)
	in = strings.Replace(in, "\r", "\n", -1)

	// Add comment to each line
	var lines []string
	for i, line := range strings.Split(in, "\n") {
		var flag string
		if i == 0 && name != "" {
			flag = name + " "
		}
		lines = append(lines, fmt.Sprintf("// %s%s", flag, line))
	}
	in = strings.Join(lines, "\n")

	// in case we have a multiline string which ends with \n, we would generate
	// empty-line-comments, like `// `. Therefore remove this line comment.
	in = strings.TrimSpace(in)
	in = strings.TrimSuffix(in, "\n//")
	return in
}

// StrInArray This function checks whether the specified string is present in an array
// of strings
func StrInArray(str string, array []string) bool {
	for _, elt := range array {
		if elt == str {
			return true
		}
	}
	return false
}

// LowercaseFirstCharacter Same as above, except lower case
func LowercaseFirstCharacter(str string) string {
	if str == "" {
		return ""
	}
	runes := []rune(str)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// UppercaseFirstCharacter Uppercase the first character in a string. This assumes UTF-8, so we have
// to be careful with unicode, don't treat it as a byte array.
func UppercaseFirstCharacter(str string) string {
	if str == "" {
		return ""
	}
	runes := []rune(str)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// Path2TypeName This converts a path, like Object/field1/nestedField into a go
// type name.
func Path2TypeName(path []string) string {
	for i, p := range path {
		path[i] = ToCamelCase(p)
	}
	return strings.Join(path, "_")
}

//FormatAndSaveFile Format code and save to file.
func FormatAndSaveFile(filename string, content []byte) error {
	var err error
	if content, err = format.Source(content); err != nil {
		return fmt.Errorf("formatAndSaveFile.source.err (%s)", err.Error())
	}
	if err = ioutil.WriteFile(filename, content, 0600); err != nil {
		return fmt.Errorf("formatAndSaveFile.err (%s)", err.Error())
	}
	return nil
}

// ExecuteTemplate 执行指定模版
func ExecuteTemplate(t *template.Template, name string, data interface{}) (string, error) {
	var buf bytes.Buffer
	var w = bufio.NewWriter(&buf)
	var err = t.ExecuteTemplate(w, name, data)
	if err != nil {
		return "", fmt.Errorf("executeTemplate.err (%s)", err.Error())
	}
	if err = w.Flush(); err != nil {
		return "", fmt.Errorf("executeTemplate.flush.err (%s)", err.Error())
	}
	return buf.String(), nil
}
