package template

import (
	"bytes"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestFuncsDocsAreUpToDate(t *testing.T) {
	data, err := ioutil.ReadFile("../../docs/Templates.md")
	if err != nil {
		t.Fatal(err)
	}

	contents := string(data)

	for funcName := range podTemplateFuncMap {
		if !strings.Contains(contents, "- `"+funcName) {
			t.Error("Function documentation missing:", funcName)
		}
	}
}

func TestExampleYamlsAreValid(t *testing.T) {
	yamlDocsPattern := regexp.MustCompile("(?sm)```yaml(.*?)```")
	templatedValue := regexp.MustCompile("{{[^{}]+}}")

	stripTemplatedBits := func(content []byte) []byte {
		var replaced []byte

		for {
			replaced = templatedValue.ReplaceAll(content, []byte("# replaced"))

			if bytes.Equal(replaced, content) {
				return replaced
			} else {
				content = replaced
			}
		}
	}

	checkYaml := func(path string, data []byte) {
		if data == nil {
			if contents, err := ioutil.ReadFile(path); err != nil {
				t.Fatal(err)
			} else {
				data = contents
			}
		}

		var source []byte
		if bytes.Contains(data, []byte("# strip-templated")) {
			source = stripTemplatedBits(data)
		} else {
			source = data
		}

		var converted map[string]interface{}
		if err := yaml.Unmarshal(source, &converted); err != nil {
			t.Error("Invalid YAML example in", path, ":", err, "\n", string(source))
		}
	}

	checkMarkdown := func(path string) {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}

		for _, match := range yamlDocsPattern.FindAll(data, -1) {
			source := yamlDocsPattern.ReplaceAll(match, []byte("$1"))
			checkYaml(path, source)
		}
	}

	filepath.Walk("../../.", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if info.Name() == "vendor" {
				return filepath.SkipDir
			}

		} else {
			switch filepath.Ext(path) {
			case ".md":
				checkMarkdown(path)
			case ".yml", ".yaml":
				checkYaml(path, nil)
			}
		}

		return nil
	})
}