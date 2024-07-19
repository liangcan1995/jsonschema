package jsonschema_test

import (
	"fmt"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"log"
	"strings"
	"testing"
)

// SchemaExt --

type objectKind struct {
	objs []jsonschema.ObjectRef
	key  string
}

func (d *objectKind) Validate(ctx *jsonschema.ValidatorContext, v any) {
	obj, ok := v.(map[string]any)
	if !ok {
		return
	}
	for _, ref := range d.objs {
		jsonschema.BuildRef(&ref, d.key, obj)
		fmt.Printf("  kind :%s ,  obj : %v", ref.Value, ref.Obj)
	}
}

// Vocab --
func ObjectKindVocab() *jsonschema.Vocabulary {
	url := "http://example.com/meta/objectKind"
	schema, err := jsonschema.UnmarshalJSON(strings.NewReader(`{
		"properties": {
			"$object-reference-kind": { "type": "string" }
		}
	}`))
	if err != nil {
		log.Fatal(err)
	}

	c := jsonschema.NewCompiler()
	if err := c.AddResource(url, schema); err != nil {
		log.Fatal(err)
	}
	sch, err := c.Compile(url)
	if err != nil {
		log.Fatal(err)
	}

	return &jsonschema.Vocabulary{
		URL:     url,
		Schema:  sch,
		Compile: compileObjectKind,
	}
}

func compileObjectKind(ctx *jsonschema.CompilerContext, obj map[string]any) (jsonschema.SchemaExt, error) {
	var keyword = "$object-reference-kind"
	// 递归遍历 JSON Schema 并提取 $ObjectRef
	objectRefs := jsonschema.ExtractObjectRefs(obj, nil, keyword)
	if objectRefs == nil || len(objectRefs) < 1 {
		return nil, nil
	}
	return &objectKind{objs: objectRefs, key: keyword}, nil
}

// Example --

func Test_vocab_ObjectKind(t *testing.T) {

	schema, err := jsonschema.UnmarshalJSON(strings.NewReader(`{
		"type": "object",
		"properties": {
			"kind": { "type": "string" },
			"fish": {
				"type": "object",
				"properties": {
					"swimmingSpeed": { 
						"type": "number"
						
					}
				},
				"required": ["swimmingSpeed"],
				"$object-reference-kind": "111111"
			},
			"dog": {
				"type": "object",
				"properties": {
					"runningSpeed": { 
						"type": "number" 
					}
				},
				"required": ["runningSpeed"],
				"$object-reference-kind": "22222"
			}
		},
		"required": ["kind"]
	}`))
	if err != nil {
		fmt.Println("xxx", err)
		log.Fatal(err)
	}
	inst, err := jsonschema.UnmarshalJSON(strings.NewReader(`{
		"kind": "fish",
		"dog":{"runningSpeed": 5},
		"fish":{"swimmingSpeed": 3}
	}`))
	if err != nil {
		log.Fatal(err)
	}
	c := jsonschema.NewCompiler()
	c.AssertVocabs()
	c.RegisterVocabulary(ObjectKindVocab())
	if err := c.AddResource("schema.json", schema); err != nil {
		log.Fatal(err)
	}
	sch, err := c.Compile("schema.json")
	if err != nil {
		log.Fatal(err)
	}

	err = sch.Validate(inst)
	fmt.Println("valid:", err == nil)
	// Output:
	// valid: false
}
