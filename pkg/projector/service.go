package projector

import (
	"github.com/aaronellington/projectl/pkg/language"
)

// NewService creates a new Projector instance
func NewService() (*Service, error) {
	languageDocker, err := language.NewDocker()
	if err != nil {
		return nil, err
	}

	languagePHP, err := language.NewPHP()
	if err != nil {
		return nil, err
	}

	languageNpm, err := language.NewNpn()
	if err != nil {
		return nil, err
	}

	languageGo, err := language.NewGo()
	if err != nil {
		return nil, err
	}

	return &Service{
		Docker: languageDocker,
		Go:     languageGo,
		PHP:    languagePHP,
		Npm:    languageNpm,
	}, nil
}

// Service is a projector
type Service struct {
	Generators  []Generator
	DistedFiles []string
	Docker      *language.Docker
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
