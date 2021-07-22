package generators

import (
	"encoding/json"
	"os"

	"github.com/aaronellington/projectl/pkg/projector"
)

// EslintConfig is the file format for the eslintrc.json
type EslintConfig struct {
	Env     map[string]bool     `json:"env"`
	Extends []string            `json:"extends"`
	Rules   map[string][]string `json:"rules"`
}

// EslintGenerator generates the .php_cs config file
type EslintGenerator struct{}

// Generate the config
func (p EslintGenerator) Generate(service *projector.Service) error {
	if !service.Npm.Enabled {
		return nil
	}

	config := &EslintConfig{
		Env: map[string]bool{
			"browser": true,
			"es2021":  true,
		},
		Extends: []string{
			"eslint:recommended",
		},
		Rules: map[string][]string{
			"comma-dangle": {"error", "always-multiline"},
			"indent":       {"error", "tab"},
			"quotes":       {"error", "single"},
			"semi":         {"error", "always"},
		},
	}

	if service.Npm.HasDependency("next") {
		config.Extends = append(config.Extends, "next")
		config.Extends = append(config.Extends, "next/core-web-vitals")
	}

	if service.Npm.HasDependency("@typescript-eslint/eslint-plugin") {
		config.Extends = append(config.Extends, "plugin:@typescript-eslint/recommended")
		config.Rules["@typescript-eslint/explicit-module-boundary-types"] = []string{"off"}
		config.Rules["@typescript-eslint/no-unused-vars"] = []string{"off"}
		config.Rules["@typescript-eslint/no-explicit-any"] = []string{"off"}
	}

	if service.Npm.HasDependency("vue") {
		config.Extends = append(config.Extends, "@vue/eslint-config-typescript/recommended")
	}

	fileBytes, _ := json.MarshalIndent(config, "", "\t")

	_ = os.WriteFile(".eslintrc.json", fileBytes, 0655)

	return nil
}
