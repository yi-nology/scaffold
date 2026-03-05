package generator

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"scaffold/internal/conf"
)

type Generator struct {
	templatePath string
	outputPath   string
	config       *conf.TemplateConfig
	variables    map[string]interface{}
	funcMap      template.FuncMap
}

type GeneratorOption func(*Generator)

func WithOutputPath(path string) GeneratorOption {
	return func(g *Generator) {
		g.outputPath = path
	}
}

func WithVariables(vars map[string]interface{}) GeneratorOption {
	return func(g *Generator) {
		g.variables = vars
	}
}

func WithFuncMap(funcMap template.FuncMap) GeneratorOption {
	return func(g *Generator) {
		g.funcMap = funcMap
	}
}

func NewGenerator(templatePath string, templateConfig *conf.TemplateConfig, opts ...GeneratorOption) *Generator {
	g := &Generator{
		templatePath: templatePath,
		config:       templateConfig,
		outputPath:   ".",
		variables:    make(map[string]interface{}),
		funcMap:      defaultFuncMap(),
	}

	for _, opt := range opts {
		opt(g)
	}

	g.applyDefaults()

	return g
}

func (g *Generator) applyDefaults() {
	for _, v := range g.config.Variables {
		if _, exists := g.variables[v.Name]; !exists {
			if v.Default != nil {
				g.variables[v.Name] = v.Default
			}
		}
	}
}

func (g *Generator) Generate() (map[string][]byte, error) {
	files := make(map[string][]byte)

	for _, fileMapping := range g.config.Files {
		if fileMapping.Exclude {
			continue
		}

		source := fileMapping.Source
		if strings.HasSuffix(source, "/*") {
			source = strings.TrimSuffix(source, "/*")
		} else if strings.HasSuffix(source, "\\*") {
			source = strings.TrimSuffix(source, "\\*")
		}

		sourcePath := filepath.Join(g.templatePath, source)

		info, err := os.Stat(sourcePath)
		if err != nil {
			return nil, fmt.Errorf("source path not found: %s", sourcePath)
		}

		if info.IsDir() {
			err = filepath.WalkDir(sourcePath, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				if d.IsDir() {
					return nil
				}

				relPath, err := filepath.Rel(sourcePath, path)
				if err != nil {
					return err
				}

				if isIgnored(relPath, g.config.Ignore) {
					return nil
				}

				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}

				targetPath := g.renderPath(fileMapping.Target, relPath)

				if fileMapping.Template != false {
					rendered, err := g.renderTemplate(string(content))
					if err != nil {
						return fmt.Errorf("failed to render %s: %w", path, err)
					}
					content = []byte(rendered)
				}

				files[targetPath] = content
				return nil
			})
		} else {
			content, err := os.ReadFile(sourcePath)
			if err != nil {
				return nil, err
			}

			targetPath := g.renderPath(fileMapping.Target, filepath.Base(sourcePath))

			if fileMapping.Template != false {
				rendered, err := g.renderTemplate(string(content))
				if err != nil {
					return nil, fmt.Errorf("failed to render %s: %w", sourcePath, err)
				}
				content = []byte(rendered)
			}

			files[targetPath] = content
		}

		if err != nil {
			return nil, err
		}
	}

	return files, nil
}

func (g *Generator) GenerateToDirectory(outputDir string) error {
	files, err := g.Generate()
	if err != nil {
		return err
	}

	for path, content := range files {
		fullPath := filepath.Join(outputDir, path)

		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return err
		}

		if err := os.WriteFile(fullPath, content, 0644); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) renderTemplate(content string) (string, error) {
	tmpl, err := template.New("content").Funcs(g.funcMap).Parse(content)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, g.variables); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (g *Generator) renderPath(baseTarget, relPath string) string {
	tmpl, err := template.New("path").Funcs(g.funcMap).Parse(baseTarget)
	if err != nil {
		return filepath.Join(baseTarget, relPath)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, g.variables); err != nil {
		return filepath.Join(baseTarget, relPath)
	}

	renderedBase := buf.String()
	return filepath.Join(renderedBase, relPath)
}

func defaultFuncMap() template.FuncMap {
	return template.FuncMap{
		"lower":     strings.ToLower,
		"upper":     strings.ToUpper,
		"title":     strings.Title,
		"trim":      strings.TrimSpace,
		"replace":   strings.ReplaceAll,
		"contains":  strings.Contains,
		"hasPrefix": strings.HasPrefix,
		"hasSuffix": strings.HasSuffix,
		"snake":     toSnakeCase,
		"camel":     toCamelCase,
		"kebab":     toKebabCase,
		"pascal":    toPascalCase,
	}
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && (r >= 'A' && r <= 'Z') {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

func toCamelCase(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})
	if len(words) == 0 {
		return s
	}
	var result strings.Builder
	result.WriteString(strings.ToLower(words[0]))
	for _, word := range words[1:] {
		if len(word) > 0 {
			result.WriteString(strings.ToUpper(string(word[0])))
			result.WriteString(strings.ToLower(word[1:]))
		}
	}
	return result.String()
}

func toKebabCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && (r >= 'A' && r <= 'Z') {
			result.WriteRune('-')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

func toPascalCase(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})
	var result strings.Builder
	for _, word := range words {
		if len(word) > 0 {
			result.WriteString(strings.ToUpper(string(word[0])))
			result.WriteString(strings.ToLower(word[1:]))
		}
	}
	return result.String()
}

func isIgnored(path string, patterns []string) bool {
	// Normalize to forward slashes for consistent matching
	normPath := filepath.ToSlash(path)
	parts := strings.Split(normPath, "/")

	for _, pattern := range patterns {
		// Check each path component for exact match
		for _, part := range parts {
			if part == pattern {
				return true
			}
			if matched, _ := filepath.Match(pattern, part); matched {
				return true
			}
		}
	}
	return false
}
