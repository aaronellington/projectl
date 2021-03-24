package projector

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
)

// GeneratorTemplated is a tempated generator
type GeneratorTemplated struct {
	Template   *template.Template
	TargetFile string
	Payload    interface{}
}

// Generate the target file
func (generator *GeneratorTemplated) Generate(service *Service) error {
	gitignoreFile, err := os.Create(generator.TargetFile)
	if err != nil {
		return fmt.Errorf("%w while creating file %s", err, generator.TargetFile)
	}

	err = generator.Template.Execute(gitignoreFile, generator.Payload)
	if err != nil {
		return fmt.Errorf("%w while executing template for %s", err, generator.TargetFile)
	}

	gitignoreFileBytes, err := os.ReadFile(generator.TargetFile)
	if err != nil {
		return fmt.Errorf("%w while reading file %s", err, generator.TargetFile)
	}

	gitignoreFileBytes = bytes.TrimSpace(gitignoreFileBytes)
	gitignoreFileBytes = append(gitignoreFileBytes, []byte("\n")...)

	err = ioutil.WriteFile(generator.TargetFile, gitignoreFileBytes, 0664)
	if err != nil {
		return fmt.Errorf("%w while writing file %s", err, generator.TargetFile)
	}

	return nil
}
