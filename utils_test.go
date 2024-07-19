package jsonschema

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_GetKeys(t *testing.T) {
	// JSON Schema 对象
	var schema map[string]any

	schemaJSON := `{
		"type": "object",
		"properties": {
			"kind": { "type": "string" },
			"fish": {
				"type": "object",
				"properties": {
					"swimmingSpeed": { "type": "number" }
				},
				"required": ["swimmingSpeed"],
				"$object-reference-kind": "11111111111"
			},
			"dog": {
				"type": "object",
				"properties": {
					"runningSpeed": { 
						"type": "number",
						"$object-reference-kind": "222222" 
					}
				},
				"required": ["runningSpeed"]
			}
		},
	    "$object-reference-kind": "333333",
		"required": ["kind"]
	}`

	// 解析 JSON Schema
	err := json.Unmarshal([]byte(schemaJSON), &schema)
	if err != nil {
		panic(err)
	}
	var keyword = "$object-reference-kind"
	// 递归遍历 JSON Schema 并提取 $ObjectRef
	objectRefs := ExtractObjectRefs(schema, nil, keyword)

	// 实例数据
	instance := `{
		"kind": "fish",
		"dog": {"runningSpeed": 5},
		"fish": {"swimmingSpeed": 3}
	}`
	// 解析实例数据
	var instanceMap map[string]any
	json.Unmarshal([]byte(instance), &instanceMap)
	// 打印结果
	for _, ref := range objectRefs {
		fmt.Printf("Path: %s, ObjectRef: %s\n", ref.Path, ref.Value)
		BuildRef(&ref, keyword, instanceMap)
		fmt.Println(ref.Obj)
	}
}
