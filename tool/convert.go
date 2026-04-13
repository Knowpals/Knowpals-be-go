package tool

import "encoding/json"

// MapToStruct map转结构体
func MapToStruct(m map[string]interface{}, outStruct interface{}) error {
	// 1. 把 map 序列化成 json 字节
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	// 2. 反序列化到结构体
	return json.Unmarshal(data, outStruct)
}
