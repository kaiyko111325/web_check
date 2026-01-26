package SheetDataToYaml

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Target struct {
	Name        string `yaml:"name"`
	URL         string `yaml:"url"`
	AllowPublic bool   `yaml:"allow_public"`
}

type TargetsConfig struct {
	Targets []Target `yaml:"targets"`
}

func SaveSheetValuesToYAML(values [][]string, filename string) error {
	config := TargetsConfig{}

	for _, row := range values {
		if len(row) < 1 {
			continue
		}
		url := row[0]
		name := ""
		if len(row) > 1 {
			name = row[1]
		}

		target := Target{
			Name:        name,
			URL:         url,
			AllowPublic: false, // 默认 false
		}
		config.Targets = append(config.Targets, target)
	}

	data, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}
