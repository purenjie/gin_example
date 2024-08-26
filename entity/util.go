package entity

import "encoding/json"

func ObjectToStr(obj interface{}) string {
	bytes, _ := json.Marshal(obj)
	return string(bytes)
}
