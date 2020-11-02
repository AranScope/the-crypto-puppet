package functions

import (
	"encoding/json"
	"reflect"
)

func min(i, j int) int {
	if i < j {
		return i
	}

	return j
}

func prettify(v interface{}) string {
	bytes, _ := json.MarshalIndent(v, "", "  ")
	return reflect.TypeOf(v).String() + "\n===============\n\n" + string(bytes) + "\n\n===============\n"
}
