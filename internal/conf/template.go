package conf

// TemplateConfig defines the structure for scaffold.yaml configuration
type TemplateConfig struct {
	Name        string   `yaml:"name" json:"name"`
	Version     string   `yaml:"version" json:"version"`
	Description string   `yaml:"description" json:"description"`
	Author      string   `yaml:"author,omitempty" json:"author,omitempty"`
	Tags        []string `yaml:"tags,omitempty" json:"tags,omitempty"`

	Repository string   `yaml:"repository,omitempty" json:"repository,omitempty"`
	Homepage   string   `yaml:"homepage,omitempty" json:"homepage,omitempty"`
	Bugs       string   `yaml:"bugs,omitempty" json:"bugs,omitempty"`
	License    string   `yaml:"license,omitempty" json:"license,omitempty"`
	Keywords   []string `yaml:"keywords,omitempty" json:"keywords,omitempty"`

	Variables []Variable    `yaml:"variables" json:"variables"`
	Files     []FileMapping `yaml:"files" json:"files"`
	Hooks     Hook          `yaml:"hooks,omitempty" json:"hooks,omitempty"`
	Ignore    []string      `yaml:"ignore,omitempty" json:"ignore,omitempty"`
}

// TemplateMeta is a simplified view for template listings
type TemplateMeta struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Author      string   `json:"author"`
	Tags        []string `json:"tags"`
	Repository  string   `json:"repository"`
	Version     string   `json:"version"`
	Homepage    string   `json:"homepage,omitempty"`
	Bugs        string   `json:"bugs,omitempty"`
	License     string   `json:"license,omitempty"`
	Keywords    []string `json:"keywords,omitempty"`
}
