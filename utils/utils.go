package utils

import "encoding/json"

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
