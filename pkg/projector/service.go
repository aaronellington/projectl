package projector

import (
	"github.com/aaronellington/projectl/pkg/language"
)

// NewService creates a new Projector instance
func NewService() (*Service, error) {
	languagePHP, err := language.NewPHP()
	if err != nil {
		return nil, err
	}

	languageNpm, err := language.NewNpm()
	if err != nil {
		return nil, err
	}

	languageGo, err := language.NewGo()
	if err != nil {
		return nil, err
	}

	return &Service{
		Go:  languageGo,
		PHP: languagePHP,
		Npm: languageNpm,
	}, nil
}

// Service is a projector
type Service struct {
	Generators  []Generator
	DistedFiles []string
	Go          *language.Go
	Npm         *language.Npm
	PHP         *language.PHP
}

// Generate your thing
func (service *Service) Generate() error {
	for _, generator := range service.Generators {

		if err := generator.Generate(service); err != nil {
			return err
		}
	}

	return nil
}

// Generator is a generator
type Generator interface {
	Generate(service *Service) error
}
