//go:generate go run ./inliner/inliner.go

package templates

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"unicode"
)

func Run(name string, tpldata interface{}) (*bytes.Buffer, error) {
	t := template.New("").Funcs(template.FuncMap{
		"ucFirst": ucFirst,
		"lcFirst": lcFirst,
		"quote":   strconv.Quote,
		"dump":    dump,
	})

	for filename, data := range data {
		_, err := t.New(filename).Parse(data)
		if err != nil {
			panic(err)
		}
	}

	buf := &bytes.Buffer{}
	err := t.Lookup(name).Execute(buf, tpldata)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func ucFirst(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func lcFirst(s string) string {
	if s == "" {
		return ""
	}

	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

func dump(val interface{}) string {
	switch val := val.(type) {
	case int:
		return strconv.Itoa(val)
	case float64:
		return fmt.Sprintf("%f", val)
	case string:
		return strconv.Quote(val)
	case bool:
		return strconv.FormatBool(val)
	case nil:
		return "nil"
	case []interface{}:
		var parts []string
		for _, part := range val {
			parts = append(parts, dump(part))
		}
		return "[]interface{}{" + strings.Join(parts, ",") + "}"
	case map[string]interface{}:
		buf := bytes.Buffer{}
		buf.WriteString("map[string]interface{}{")
		for key, data := range val {
			buf.WriteString(strconv.Quote(key))
			buf.WriteString(":")
			buf.WriteString(dump(data))
		}
		buf.WriteString("}")
		return buf.String()
	default:
		panic(fmt.Errorf("unsupported type %T", val))
	}
}