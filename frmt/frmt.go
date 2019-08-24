package frmt

import (
	"bytes"
	"fmt"
	"text/template"
)

// Render formats a string template and outputs console ready text
func Render(in string, data interface{}) []byte {
	tpl, err := template.New("").Funcs(funcMap).Parse(in)
	if err != nil {
		return []byte(fmt.Sprintf("%v", err))
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	if err != nil {
		return []byte(fmt.Sprintf("%v", err))
	}

	return buf.Bytes()
}
