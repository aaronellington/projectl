package generators

import (
	"os"

	"github.com/aaronellington/projectl/pkg/projector"
)

// GithubWorkflow generator
type GithubWorkflow struct{}

// Generate the config file
func (githubWorkflow *GithubWorkflow) Generate(service *projector.Service) error {
	if err := os.MkdirAll(".github/workflows", 0775); err != nil {
		return err
	}

	workflowFile, err := os.Create(".github/workflows/main.yml")
	if err != nil {
		return err
	}

	workflowFile.WriteString(`name: Main

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
`)

	if service.Npm.Enabled {
		workflowFile.WriteString(`
      - name: Set up Node
        uses: actions/setup-node@v1
        with:
          node-version: 14
`)
	}

	if service.Go.Enabled {
		workflowFile.WriteString(`
      - name: Set up Go
        uses: actions/setup-go@v1
        with:
          go-version: ` + service.Go.Version() + `
`)
	}

	if service.PHP.Enabled {
		workflowFile.WriteString(`
      - name: Set up PHP
        uses: shivammathur/setup-php@v2
        with:
          php-version: '8.0'
          tools: composer:v2
`)
	}

	workflowFile.WriteString(`
      - name: Check out code
        uses: actions/checkout@v2
`)
	workflowFile.WriteString(`
      - name: Build
        run: make lint test build git-change-check
`)

	return nil
}
