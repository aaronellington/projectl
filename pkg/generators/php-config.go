package generators

import (
	"os"

	"github.com/aaronellington/projectl/pkg/projector"

	// For embed
	_ "embed"
)

//go:embed php/php_cs.php
var phpCSFixerConfigFile []byte

//go:embed php/phpcs.xml
var phpCodeSnifferConfigFile []byte

// PHPConfig generates the .php_cs config file
type PHPConfig struct{}

// Generate the config
func (p PHPConfig) Generate(service *projector.Service) error {
	if !service.PHP.Enabled {
		return nil
	}

	if err := os.WriteFile(".php_cs", phpCSFixerConfigFile, 0655); err != nil {
		return err
	}

	if err := os.WriteFile(".phpcs.xml", phpCodeSnifferConfigFile, 0655); err != nil {
		return err
	}

	return nil
}
