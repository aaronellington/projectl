package generators

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/aaronellington/projectl/pkg/configuration"
	"github.com/aaronellington/projectl/pkg/projector"
)

const makefileTemplate = `.PHONY:{{ range .Targets }} {{ .Name }}{{ end }}

SHELL=/bin/bash -o pipefail

{{ range $Key, $Value := .Variables }}{{ $Key }} := {{ $Value }}
{{ end }}
{{ range .Targets }}{{ .Name }}:{{ range .PreTargets }} {{ . }}{{ end }}{{ if .Comment }} ## {{ .Comment}}{{ end }}
{{ range .Commands }}	{{ . }}
{{ end }}
{{ end }}
`

// TemplatePayloadMakefile template payload
type TemplatePayloadMakefile struct {
	Variables map[string]string
	Targets   []*TemplateMakefileTarget
}

// TemplateMakefileTarget is a target of Makefile
type TemplateMakefileTarget struct {
	Name       string
	Comment    string
	PreTargets []string
	Commands   []string
}

// NewMakefile generator
func NewMakefile(service *projector.Service, config *configuration.Config) *projector.GeneratorTemplated {
	template := template.Must(template.New("makefile").Parse(makefileTemplate))

	return &projector.GeneratorTemplated{
		TargetFile: "Makefile",
		Template:   template,
		Payload:    getMakefilePayload(service, config),
	}
}

func getMakefilePayload(service *projector.Service, config *configuration.Config) TemplatePayloadMakefile {
	payload := &TemplatePayloadMakefile{
		Variables: make(map[string]string),
	}

	payload.Variables[".DEFAULT_GOAL"] = "help"

	if service.Go.Enabled {
		payload.Variables["GO_PATH"] = "$(shell go env GOPATH 2> /dev/null)"
		payload.Variables["PATH"] = "$(GO_PATH)/bin:$(PATH)"
	}
	if service.PHP.Enabled {
		payload.Variables["COMPOSER_BIN"] = "$(shell composer config bin-dir 2> /dev/null)"
	}

	addHelpTarget(service, payload)
	addDockerTargets(service, config, payload)
	addBuildTargets(service, payload)
	addLintTargets(service, payload)
	addTestTargets(service, payload)
	addWatchTargets(service, config, payload)
	addCleanTargets(service, payload)
	addCopyConfigTarget(service, payload)
	addPipelineTargets(service, payload)
	addAnsibleTargets(payload)

	// TODO: add target for fix
	// TODO: add target for fix-php
	// TODO: add target for fix-npm
	// TODO: add target for fix-go

	return *payload
}

func addDockerTargets(service *projector.Service, config *configuration.Config, payload *TemplatePayloadMakefile) {
	if !service.Docker.Enabled {
		return
	}

	targetDocker := &TemplateMakefileTarget{
		Name: "docker",
		Commands: []string{
			"docker build -t " + config.DockerName + ":latest .",
		},
	}

	payload.Targets = append(payload.Targets, targetDocker)
}

func addPipelineTargets(service *projector.Service, payload *TemplatePayloadMakefile) {
	targetPostLint := &TemplateMakefileTarget{
		Name: "git-change-check",
		Commands: []string{
			"@git diff --exit-code --quiet || (echo 'There should not be any changes at this point' && git status && exit 1;)",
		},
	}
	payload.Targets = append(payload.Targets, targetPostLint)
}

func addCopyConfigTarget(service *projector.Service, payload *TemplatePayloadMakefile) {
	commands := []string{}
	for _, distedFile := range service.DistedFiles {
		if _, err := os.Open("." + distedFile + ".dist"); err != nil {
			continue
		}

		commands = append(commands, fmt.Sprintf("[ -f %s ] || cp %s.dist %s", distedFile, distedFile, distedFile))
	}

	targetCopyConfig := &TemplateMakefileTarget{
		Name:     "copy-config",
		Comment:  "Copy missing config files into place",
		Commands: commands,
	}
	payload.Targets = append(payload.Targets, targetCopyConfig)
}

func addHelpTarget(service *projector.Service, payload *TemplatePayloadMakefile) {
	targetHelp := &TemplateMakefileTarget{
		Name:    "help",
		Comment: "Display general help about this command",
		Commands: []string{
			"@echo 'Makefile targets:'",
			"@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' Makefile \\",
			"| sed -n 's/^\\(.*\\): \\(.*\\)##\\(.*\\)/    \\1 :: \\3/p' \\",
			"| column -t -c 1  -s '::'",
		},
	}
	payload.Targets = append(payload.Targets, targetHelp)
}

func addCleanTargets(service *projector.Service, payload *TemplatePayloadMakefile) {
	cleanArguments := []string{}
	for _, distedFile := range service.DistedFiles {
		if _, err := os.Open("." + distedFile + ".dist"); err != nil {
			continue
		}

		cleanArguments = append(cleanArguments, fmt.Sprintf(" --exclude='!%s'", distedFile))
	}

	targetClean := &TemplateMakefileTarget{
		Name:    "clean",
		Comment: "Remove files listed in .gitignore (possibly with some exceptions)",
		Commands: []string{
			"@git init 2> /dev/null",
			"git clean -Xdff" + strings.Join(cleanArguments, ""),
		},
	}
	payload.Targets = append(payload.Targets, targetClean)

	targetCleanFull := &TemplateMakefileTarget{
		Name: "clean-full",
		Commands: []string{
			"@git init 2> /dev/null",
			"git clean -Xdff",
		},
	}
	payload.Targets = append(payload.Targets, targetCleanFull)
}

func addBuildTargets(service *projector.Service, payload *TemplatePayloadMakefile) {
	targetBuild := &TemplateMakefileTarget{
		Name:    "build",
		Comment: "Build the application",
	}
	payload.Targets = append(payload.Targets, targetBuild)

	if service.Npm.Enabled {
		targetBuildNpm := &TemplateMakefileTarget{
			Name: "build-npm",
			Commands: []string{
				"npm install",
				"npm run build",
			},
		}
		targetBuild.PreTargets = append(targetBuild.PreTargets, targetBuildNpm.Name)
		payload.Targets = append(payload.Targets, targetBuildNpm)
	}

	if service.PHP.Enabled {
		targetBuildPHPProd := &TemplateMakefileTarget{
			Name: "build-php-prod",

			Commands: []string{
				"composer install --no-dev --optimize-autoloader --classmap-authoritative --no-progress --no-interaction",
			},
		}
		if service.PHP.IsSymfony3() {
			targetBuildPHPProd.Commands = []string{
				"SYMFONY_ENV=prod composer install --no-dev --optimize-autoloader --classmap-authoritative --no-progress --no-interaction",
				"rsync -a --exclude='web/app_*.php' --exclude='var/cache' --exclude='/vendor/**/.git' app var web bin src vendor sass node_modules js composer.json httpsdocs",
			}
		}
		payload.Targets = append(payload.Targets, targetBuildPHPProd)
		targetBuild.PreTargets = append(targetBuild.PreTargets, targetBuildPHPProd.Name)

		targetBuildPHPTest := &TemplateMakefileTarget{
			Name: "build-php-test",
			Commands: []string{
				"composer install --no-progress --no-interaction",
			},
		}
		payload.Targets = append(payload.Targets, targetBuildPHPTest)
	}

	if service.Go.Enabled {
		buildCommands := []string{
			"@go generate",
		}

		for path, name := range service.Go.Targets {
			buildCommands = append(buildCommands,
				"go build -ldflags='-s -w' -o $(CURDIR)/var/"+name+" "+path,
				"@ln -sf $(CURDIR)/var/"+name+" $(GO_PATH)/bin/"+name,
			)
		}

		targetBuildGo := &TemplateMakefileTarget{
			Name:     "build-go",
			Commands: buildCommands,
		}
		targetBuild.PreTargets = append(targetBuild.PreTargets, targetBuildGo.Name)
		payload.Targets = append(payload.Targets, targetBuildGo)
	}
}

func addTestTargets(service *projector.Service, payload *TemplatePayloadMakefile) {
	targetTest := &TemplateMakefileTarget{
		Name:    "test",
		Comment: "Test the application",
	}
	payload.Targets = append(payload.Targets, targetTest)

	if service.Npm.Enabled {
		targetTestNpm := &TemplateMakefileTarget{
			Name: "test-npm",
			Commands: []string{
				"npm install",
				"npm run test",
			},
		}
		targetTest.PreTargets = append(targetTest.PreTargets, targetTestNpm.Name)
		payload.Targets = append(payload.Targets, targetTestNpm)
	}

	if service.PHP.Enabled {
		targetTestPHP := &TemplateMakefileTarget{
			Name:       "test-php",
			PreTargets: []string{"build-php-test"},
			Commands: []string{
				"$(COMPOSER_BIN)/phpunit src",
			},
		}
		targetTest.PreTargets = append(targetTest.PreTargets, targetTestPHP.Name)
		payload.Targets = append(payload.Targets, targetTestPHP)
	}

	if service.Go.Enabled {
		targetTestGo := &TemplateMakefileTarget{
			Name: "test-go",
			Commands: []string{
				"@mkdir -p var/",
				"@go test -race -cover -coverprofile  var/coverage.txt ./...",
				"@go tool cover -func var/coverage.txt | awk '/^total/{print $$1 \" \" $$3}'",
			},
		}
		targetTest.PreTargets = append(targetTest.PreTargets, targetTestGo.Name)
		payload.Targets = append(payload.Targets, targetTestGo)
	}
}

func addWatchTargets(service *projector.Service, config *configuration.Config, payload *TemplatePayloadMakefile) {
	if service.Npm.Enabled && service.Npm.HasScript("watch") {
		targetTestNpm := &TemplateMakefileTarget{
			Name: "watch-npm",
			Commands: []string{
				"clear",
				"npm run watch",
			},
		}
		payload.Targets = append(payload.Targets, targetTestNpm)
	}

	if service.Go.Enabled && config.GoHTTP {
		// TODO: does not support multiple go targets
		targetTestGo := &TemplateMakefileTarget{
			Name: "watch-go",
			Commands: []string{
				"@cd ; go get github.com/codegangsta/gin",
				"clear",
				"gin --all --immediate --path . --build . --bin var/gin run",
			},
		}
		payload.Targets = append(payload.Targets, targetTestGo)
	}
}

func addLintTargets(service *projector.Service, payload *TemplatePayloadMakefile) {
	targetLint := &TemplateMakefileTarget{
		Name:    "lint",
		Comment: "Lint the application",
	}
	payload.Targets = append(payload.Targets, targetLint)

	if service.Npm.Enabled {
		targetLanguage := &TemplateMakefileTarget{
			Name: "lint-npm",
			Commands: []string{
				"npm install",
				"npm run lint",
			},
		}
		targetLint.PreTargets = append(targetLint.PreTargets, targetLanguage.Name)
		payload.Targets = append(payload.Targets, targetLanguage)
	}

	if service.PHP.Enabled {
		targetLintPHP := &TemplateMakefileTarget{
			Name:       "lint-php",
			PreTargets: []string{"build-php-test"},
			Commands: []string{
				"$(COMPOSER_BIN)/php-cs-fixer fix",
				"$(COMPOSER_BIN)/phpcs",
				"$(COMPOSER_BIN)/phpstan analyse src --level=max",
			},
		}
		targetLint.PreTargets = append(targetLint.PreTargets, targetLintPHP.Name)
		payload.Targets = append(payload.Targets, targetLintPHP)
	}

	if service.Go.Enabled {
		targetLintGo := &TemplateMakefileTarget{
			Name: "lint-go",
			Commands: []string{
				"@cd ; go get golang.org/x/lint/golint",
				"@cd ; go get golang.org/x/tools/cmd/goimports",
				"go get -d ./...",
				"go mod tidy",
				"gofmt -s -w .",
				"go vet ./...",
				"golint -set_exit_status=1 ./...",
				"goimports -w .",
			},
		}
		targetLint.PreTargets = append(targetLint.PreTargets, targetLintGo.Name)
		payload.Targets = append(payload.Targets, targetLintGo)
	}
}
func addAnsibleTargets(payload *TemplatePayloadMakefile) {
	files, err := ioutil.ReadDir("ansible/playbooks/")
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Fatal(err)
		}
	}
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".yml") {
			continue
		}
		payload.Targets = append(payload.Targets, &TemplateMakefileTarget{
			Name:       "ansible-" + strings.TrimSuffix(file.Name(), ".yml"),
			PreTargets: []string{"git-change-check"},
			Commands: []string{
				"clear",
				"time ansible-playbook ansible/playbooks/" + file.Name(),
			},
		})
	}

}
