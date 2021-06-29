package generators

import (
	"os"
	"text/template"

	"github.com/aaronellington/projectl/pkg/configuration"
	"github.com/aaronellington/projectl/pkg/projector"
)

const gitignoreTemplate = `{{ range .Sections }}{{ if .Values }}# {{ .Name }}{{ range .Values }}
{{ . }}{{ end }}

{{ end }}{{ end }}
`

// TemplatePayloadGitignore template payload
type TemplatePayloadGitignore struct {
	Sections []TemplateGitignoreSection
}

// TemplateGitignoreSection is a section of gitignore
type TemplateGitignoreSection struct {
	Name   string
	Values []string
}

// NewGitignore generator
func NewGitignore(service *projector.Service, config *configuration.Config) *projector.GeneratorTemplated {
	template := template.Must(template.New("gitignore").Parse(gitignoreTemplate))

	return &projector.GeneratorTemplated{
		TargetFile: ".gitignore",
		Template:   template,
		Payload:    getGitignorePayload(service, config),
	}
}

func getGitignorePayload(service *projector.Service, config *configuration.Config) TemplatePayloadGitignore {
	distedFiles := []string{}
	for _, distedFile := range service.DistedFiles {
		if _, err := os.Open("." + distedFile + ".dist"); err != nil {
			continue
		}

		distedFiles = append(distedFiles, distedFile)
	}

	payload := TemplatePayloadGitignore{
		Sections: []TemplateGitignoreSection{
			{
				Name: "System Files",
				Values: []string{
					"/.vscode/",
					"/.idea/",
					".DS_Store",
				},
			},
			{
				Name: "Temporary Files",
				Values: []string{
					"/var/",
				},
			},
			{
				Name:   "Disted Files",
				Values: distedFiles,
			},
			{
				Name: "Environment Files",
				Values: []string{
					"/.env.local",
					"/.env.*.local",
				},
			},
		},
	}

	// Add PHP Values
	if service.PHP.Enabled {
		payload.Sections = append(payload.Sections, TemplateGitignoreSection{
			Name: "PHP Files",
			Values: []string{
				"/vendor/",
				".phpunit.result.cache",
				".php_cs.cache",
				".phpcs-cache",
			},
		})
	}

	// Add NPM Values
	if service.Npm.Enabled {
		npmValues := []string{
			"/node_modules/",
			"npm-debug.log",
		}

		if service.Go.Enabled {
			npmValues = append(npmValues, "/resources/dist/")
		}
		if service.PHP.Enabled {
			npmValues = append(npmValues, "/public/build/")
		}

		payload.Sections = append(payload.Sections, TemplateGitignoreSection{
			Name:   "NPM Files",
			Values: npmValues,
		})

	}

	// Add Go Values
	if service.Go.Enabled {
		goGitignoreValues := []string{
			"__debug_bin",
			"debug.test",
		}

		payload.Sections = append(payload.Sections, TemplateGitignoreSection{
			Name:   "Go Files",
			Values: goGitignoreValues,
		})
	}

	// Add Project Specific Files
	payload.Sections = append(payload.Sections, TemplateGitignoreSection{
		Name:   "Project Specific Files",
		Values: config.Gitignore,
	})

	return payload
}
