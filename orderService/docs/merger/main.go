package main

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fileContents := make([]string, 0)
	if err := filepath.WalkDir("./docs/api", walk(&fileContents)); err != nil {
		log.Fatal(err)
	}

	err := getSwaggerYAML(fileContents)
	if err != nil {
		log.Fatal(err)
	}

	err = getSwaggerJSON(fileContents)
	if err != nil {
		log.Fatal(err)
	}
}

func getSwaggerYAML(fileContents []string) error {
	resultYAML, err := mergeYamlValues(fileContents)
	if err != nil {
		return err
	}

	if err = os.WriteFile("./docs/api.gen.yaml", []byte(resultYAML), 0o644); err != nil {
		return err
	}

	return nil
}

func getSwaggerJSON(fileContents []string) error {
	resultYAML, err := mergeYamlValues(fileContents)
	if err != nil {
		log.Fatal(err)
	}

	var body interface{}
	if err := yaml.Unmarshal([]byte(resultYAML), &body); err != nil {
		return err
	}
	body = convertToJSON(body)

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	if err = os.WriteFile("./docs/dist/swagger.json", b, 0o644); err != nil {
		return err
	}

	return nil
}

func walk(contents *[]string) func(path string, d fs.DirEntry, err error) error {
	return func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Exclude oapi codegen config file
		if entry.Name() == "entities.cfg.yaml" || entry.Name() == "server.cfg.yaml" {
			return nil
		}

		fileName := strings.ToLower(path)

		if strings.HasSuffix(fileName, "yaml") || strings.HasSuffix(fileName, "yml") {
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			*contents = append(*contents, string(b))
		}

		return nil
	}
}

func mergeYamlValues(values []string) (string, error) {
	var result map[any]any

	var bs []byte

	for _, value := range values {
		var override map[any]any

		bs = []byte(value)

		if err := yaml.Unmarshal(bs, &override); err != nil {
			return "", err
		}

		// check if is nil. This will only happen for the first value
		if result == nil {
			result = override
		} else {
			result = mergeMaps(result, override)
		}
	}

	bs, err := yaml.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}

func mergeMaps(a, b map[any]any) map[any]any {
	out := make(map[any]any, len(a))
	for k, v := range a {
		out[k] = v
	}

	for k, v := range b {
		if v, ok := v.(map[any]any); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[any]any); ok {
					out[k] = mergeMaps(bv, v)
					continue
				}
			}
		}

		out[k] = v
	}

	return out
}

func convertToJSON(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convertToJSON(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convertToJSON(v)
		}
	}
	return i
}
