package main

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"os"
)

func CreateTargetDirectory(year int, dayName string) (string, error) {
	targetDirectory := fmt.Sprintf("%d/%s", year, dayName)
	err := os.MkdirAll(targetDirectory, 0777)
	switch {
	case errors.Is(err, os.ErrExist):
	case err != nil:
		return "", fmt.Errorf("make directory %s: %w", targetDirectory, err)
	}

	return targetDirectory, nil
}

type Readme struct {
	data string
}

func NewReadme(data string) *Readme {
	return &Readme{data: data}
}

func (r *Readme) CreateFile(dir string) error {
	err := os.WriteFile(fmt.Sprintf("%s/Readme.md", dir), []byte(r.data), 0777)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

//go:embed templates
var templates embed.FS

type Go struct {
	dayName string
	year    int
}

func NewGo(dayName string, year int) *Go {
	return &Go{
		dayName: dayName,
		year:    year,
	}
}

func (g *Go) CreateFile(dir string) error {
	filledTemplate, err := g.fillTemplate()
	if err != nil {
		return fmt.Errorf("fill template: %w", err)
	}

	err = os.WriteFile(fmt.Sprintf("%s/%s.go", dir, g.dayName), filledTemplate, 0777)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

func (g *Go) fillTemplate() ([]byte, error) {
	tmpl, err := template.ParseFS(templates, "templates/solve.go.tmpl")
	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}

	out := bytes.NewBuffer(nil)

	err = tmpl.Execute(out, map[string]interface{}{"Day": g.dayName, "Year": g.year})
	if err != nil {
		return nil, fmt.Errorf("execute template: %w", err)
	}

	return out.Bytes(), nil
}

type InputType int

const (
	TestInputType = InputType(iota)
	ProblemInputType
)

func (i InputType) Filename() string {
	switch i {
	case TestInputType:
		return "test.txt"
	case ProblemInputType:
		return "input.txt"
	default:
		return "invalid-input-type.txt"
	}
}

type Input struct {
	inputType InputType
	data      string
}

func NewInput(inputType InputType, data string) *Input {
	return &Input{inputType: inputType, data: data}
}

func (i *Input) CreateFile(dir string) error {
	if i.data[len(i.data)-1] == '\n' {
		i.data = i.data[:len(i.data)-1]
	}

	err := os.WriteFile(fmt.Sprintf("%s/%s", dir, i.inputType.Filename()), []byte(i.data), 0777)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
