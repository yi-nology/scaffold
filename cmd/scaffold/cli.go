package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"scaffold/internal/app/generator"
	tmplservice "scaffold/internal/app/template"
	"scaffold/internal/conf"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7C3AED")).
			MarginBottom(1)

	questionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#60A5FA")).
			MarginRight(1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981"))

	optionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EF4444"))

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7C3AED")).
			Padding(1, 2).
			Margin(1, 0)
)

type step int

const (
	stepTemplateSelect step = iota
	stepVariables
	stepOutputDir
	stepGenerating
	stepDone
)

type cliModel struct {
	service      *tmplservice.Service
	templates    []conf.TemplateMeta
	selectedTmpl int
	step         step

	variables  []conf.Variable
	varIndex   int
	varInputs  map[string]interface{}
	textInputs []textinput.Model
	inputIndex int

	outputDir   string
	outputInput textinput.Model

	selectedTemplate *tmplservice.TemplateSource
	err              error
}

func runBubbleTeaUI(service *tmplservice.Service, projectName string) error {
	templates := service.ListTemplates()

	initialInputs := make(map[string]interface{})
	if projectName != "" {
		initialInputs["project_name"] = projectName
	}

	m := cliModel{
		service:   service,
		templates: templates,
		step:      stepTemplateSelect,
		varInputs: initialInputs,
		outputDir: ".",
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	if fm, ok := finalModel.(cliModel); ok {
		return fm.err
	}
	return nil
}

func (m cliModel) Init() tea.Cmd {
	return nil
}

func (m cliModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			return m.handleUp(), nil
		case "down", "j":
			return m.handleDown(), nil
		case "enter":
			return m.handleEnter()
		case "tab":
			return m.handleTab(), nil
		}
	}

	if m.step == stepVariables && len(m.textInputs) > 0 {
		var cmd tea.Cmd
		ti := &m.textInputs[m.inputIndex]
		*ti, cmd = ti.Update(msg)
		return m, cmd
	}

	if m.step == stepOutputDir {
		var cmd tea.Cmd
		m.outputInput, cmd = m.outputInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *cliModel) handleUp() cliModel {
	switch m.step {
	case stepTemplateSelect:
		if m.selectedTmpl > 0 {
			m.selectedTmpl--
		}
	case stepVariables:
		if m.inputIndex > 0 {
			m.inputIndex--
		}
	}
	return *m
}

func (m *cliModel) handleDown() cliModel {
	switch m.step {
	case stepTemplateSelect:
		if m.selectedTmpl < len(m.templates)-1 {
			m.selectedTmpl++
		}
	case stepVariables:
		if m.inputIndex < len(m.textInputs)-1 {
			m.inputIndex++
		}
	}
	return *m
}

func (m *cliModel) handleTab() cliModel {
	if m.step == stepVariables && len(m.textInputs) > 0 {
		if m.inputIndex < len(m.textInputs)-1 {
			m.inputIndex++
		} else {
			m.inputIndex = 0
		}
	}
	return *m
}

func (m *cliModel) handleEnter() (tea.Model, tea.Cmd) {
	switch m.step {
	case stepTemplateSelect:
		if len(m.templates) == 0 {
			m.err = fmt.Errorf("no templates available")
			return m, tea.Quit
		}

		tmpl := m.templates[m.selectedTmpl]
		source, err := m.service.GetTemplate(tmpl.ID)
		if err != nil {
			m.err = err
			return m, tea.Quit
		}

		m.selectedTemplate = source
		m.variables = source.Config.Variables
		m.step = stepVariables

		m.textInputs = make([]textinput.Model, 0)
		for _, v := range m.variables {
			ti := textinput.New()
			ti.Placeholder = fmt.Sprintf("%v", v.Default)
			ti.EchoMode = textinput.EchoNormal
			if v.Type == conf.VariableTypeBoolean {
				ti.EchoMode = textinput.EchoNone
			}

			if val, exists := m.varInputs[v.Name]; exists {
				ti.SetValue(fmt.Sprintf("%v", val))
			}

			m.textInputs = append(m.textInputs, ti)
		}
		if len(m.textInputs) > 0 {
			m.textInputs[0].Focus()
		}

	case stepVariables:
		for i, v := range m.variables {
			val := m.textInputs[i].Value()
			if val == "" {
				m.varInputs[v.Name] = v.Default
			} else {
				switch v.Type {
				case conf.VariableTypeBoolean:
					m.varInputs[v.Name] = strings.ToLower(val) == "true" || val == "yes" || val == "y"
				default:
					m.varInputs[v.Name] = val
				}
			}
		}

		m.step = stepOutputDir
		m.outputInput = textinput.New()
		m.outputInput.SetValue(m.outputDir)
		m.outputInput.Focus()

	case stepOutputDir:
		m.outputDir = m.outputInput.Value()
		if m.outputDir == "" {
			m.outputDir = "."
		}
		m.step = stepGenerating
		return m, m.generateProject()

	case stepDone:
		return m, tea.Quit
	}

	return m, nil
}

func (m *cliModel) generateProject() tea.Cmd {
	return func() tea.Msg {
		gen := generator.NewGenerator(
			m.selectedTemplate.LocalPath,
			m.selectedTemplate.Config,
			generator.WithVariables(m.varInputs),
		)

		outputPath := m.outputDir
		if projectName, ok := m.varInputs["project_name"].(string); ok && projectName != "" {
			outputPath = filepath.Join(m.outputDir, projectName)
		}

		if err := os.MkdirAll(outputPath, 0755); err != nil {
			m.err = err
			m.step = stepDone
			return nil
		}

		if err := gen.GenerateToDirectory(outputPath); err != nil {
			m.err = err
			m.step = stepDone
			return nil
		}

		m.step = stepDone
		return nil
	}
}

func (m cliModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Scaffold - Project Generator"))
	b.WriteString("\n\n")

	switch m.step {
	case stepTemplateSelect:
		b.WriteString("Select a template:\n\n")
		for i, t := range m.templates {
			cursor := " "
			if m.selectedTmpl == i {
				cursor = selectedStyle.Render(">")
			}
			b.WriteString(fmt.Sprintf("%s %s - %s\n", cursor, t.Name, t.Description))
		}

	case stepVariables:
		if len(m.variables) == 0 {
			b.WriteString("No variables to configure.\n")
		} else {
			b.WriteString("Configure your project:\n\n")
			for i, v := range m.variables {
				cursor := " "
				if m.inputIndex == i {
					cursor = selectedStyle.Render(">")
				}

				prompt := v.Prompt
				if prompt == "" {
					prompt = v.Name
				}

				var input string
				if m.inputIndex == i {
					input = m.textInputs[i].View()
				} else {
					val := m.textInputs[i].Value()
					if val == "" {
						val = fmt.Sprintf("%v", v.Default)
					}
					input = optionStyle.Render(val)
				}

				b.WriteString(fmt.Sprintf("%s %s: %s\n", cursor, questionStyle.Render(prompt), input))
			}
			b.WriteString("\n")
			b.WriteString(optionStyle.Render("Press Enter to continue, Tab to switch fields"))
		}

	case stepOutputDir:
		b.WriteString("Output directory:\n\n")
		b.WriteString(fmt.Sprintf("%s %s\n", questionStyle.Render("Path:"), m.outputInput.View()))
		b.WriteString("\n")
		b.WriteString(optionStyle.Render("Press Enter to generate project"))

	case stepGenerating:
		b.WriteString("Generating project...\n")

	case stepDone:
		if m.err != nil {
			b.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
		} else {
			outputPath := m.outputDir
			if projectName, ok := m.varInputs["project_name"].(string); ok && projectName != "" {
				outputPath = filepath.Join(m.outputDir, projectName)
			}
			absPath, _ := filepath.Abs(outputPath)
			b.WriteString(selectedStyle.Render("Project generated successfully!"))
			b.WriteString(fmt.Sprintf("\n\nOutput: %s", absPath))
			b.WriteString("\n\nPress Enter to exit")
		}
	}

	return boxStyle.Render(b.String())
}
