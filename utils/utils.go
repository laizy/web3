package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func Ensure(err error) {
	if err != nil {
		panic(err)
	}
}

func JsonString(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	Ensure(err)

	return string(b)
}

func JsonStr(v interface{}) string {
	b, err := json.Marshal(v)
	Ensure(err)

	return string(b)
}

func EnsureTrue(b bool) {
	if !b {
		panic("must be true")
	}
}

func LoadJsonFile(file string, val interface{}) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(content, val); err != nil {
		if syntaxerr, ok := err.(*json.SyntaxError); ok {
			line := findLine(content, syntaxerr.Offset)
			return fmt.Errorf("JSON syntax error at %v:%v: %v", file, line, err)
		}
		return fmt.Errorf("JSON unmarshal error in %v: %v", file, err)
	}
	return nil
}

// findLine returns the line number for the given offset into data.
func findLine(data []byte, offset int64) (line int) {
	line = 1
	for i, r := range string(data) {
		if int64(i) >= offset {
			return
		}
		if r == '\n' {
			line++
		}
	}
	return
}
