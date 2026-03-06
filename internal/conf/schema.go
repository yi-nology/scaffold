package conf

type VariableType string

const (
	VariableTypeString  VariableType = "string"
	VariableTypeBoolean VariableType = "boolean"
	VariableTypeEnum    VariableType = "enum"
	VariableTypeNumber  VariableType = "number"
)

type Variable struct {
	Name        string       `yaml:"name" json:"name"`
	Type        VariableType `yaml:"type" json:"type"`
	Default     interface{}  `yaml:"default" json:"default"`
	Prompt      string       `yaml:"prompt" json:"prompt"`
	Options     []string     `yaml:"options,omitempty" json:"options,omitempty"`
	Required    bool         `yaml:"required" json:"required"`
	Validation  string       `yaml:"validation,omitempty" json:"validation,omitempty"`
	Description string       `yaml:"description,omitempty" json:"description,omitempty"`
	Group       string       `yaml:"group,omitempty" json:"group,omitempty"`
}

type FileMapping struct {
	Source   string `yaml:"source" json:"source"`
	Target   string `yaml:"target" json:"target"`
	Exclude  bool   `yaml:"exclude,omitempty" json:"exclude,omitempty"`
	Template bool   `yaml:"template,omitempty" json:"template,omitempty"`
}

type Hook struct {
	Pre  string `yaml:"pre,omitempty" json:"pre,omitempty"`
	Post string `yaml:"post,omitempty" json:"post,omitempty"`
}
