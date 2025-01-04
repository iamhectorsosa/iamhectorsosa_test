package templates

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Heading struct {
	Name    string        `yaml:"name"`
	Bio     string        `yaml:"bio"`
	About   string        `yaml:"about"`
	Repo    string        `yaml:"repo"`
	URL     string        `yaml:"url"`
	Contact []ContactItem `yaml:"contact"`
}

type ContactItem struct {
	Type  string `yaml:"type"`
	Href  string `yaml:"href"`
	Label string `yaml:"label"`
}

type Section struct {
	Title string        `yaml:"title"`
	Items []SectionItem `yaml:"items"`
}

type SectionItem struct {
	Title       string  `yaml:"title"`
	StartDate   string  `yaml:"start_date"`
	EndDate     *string `yaml:"end_date"`
	Href        string  `yaml:"href,omitempty"`
	Description string  `yaml:"description,omitempty"`
}

type Config struct {
	Heading  Heading   `yaml:"heading"`
	Sections []Section `yaml:"sections"`
}

func ParseYaml(filePath string) (Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
