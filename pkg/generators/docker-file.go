package generators

import (
	"fmt"
	"os"

	"github.com/aaronellington/projectl/pkg/projector"
)

// Dockerfile generator
type Dockerfile struct {
	Port   int
	Target string
}

// Generate the config file
func (dockerfile *Dockerfile) Generate(service *projector.Service) error {
	file, err := os.Create("Dockerfile")
	if err != nil {
		return err
	}

	if service.Npm.Enabled {
		file.WriteString(`FROM node:16-buster as nodeBuilder
WORKDIR /build-staging
COPY . .
RUN make clean-full
RUN make lint-npm test-npm build-npm

`)
	}

	if service.Go.Enabled {
		dockerfile.modeGo(service, file)
	} else if service.PHP.Enabled {
		dockerfile.modePHP(service, file)
	}

	if dockerfile.Port != 0 {
		file.WriteString(fmt.Sprintf("EXPOSE %d\n", dockerfile.Port))
	}

	return nil
}

func (dockerfile *Dockerfile) modeGo(service *projector.Service, file *os.File) {
	targetBin := dockerfile.Target

	if targetBin == "" {
		for _, bin := range service.Go.Targets {
			targetBin = bin
			break
		}
	}

	file.WriteString(`FROM golang:` + service.Go.Version() + `-buster as goBuilder
WORKDIR /build-staging
COPY . .
RUN make clean-full
`)

	if service.Npm.Enabled {
		file.WriteString("COPY --from=nodeBuilder /build-staging/resources/dist/ /build-staging/resources/dist/\n")
	}
	file.WriteString(`RUN make lint-go test-go build-go

FROM debian:buster
RUN apt-get update
RUN apt-get install -y ca-certificates
WORKDIR /app
`)

	file.WriteString(`COPY --from=goBuilder /build-staging/var/` + targetBin + ` ./` + targetBin + `
CMD ["./` + targetBin + `"]
`)
}

func (dockerfile *Dockerfile) modePHP(service *projector.Service, file *os.File) {
	file.WriteString(`FROM aaronellington/php-fpm-webserver:latest
COPY . .
RUN make clean-full
`)

	if service.Npm.Enabled {
		file.WriteString("COPY --from=nodeBuilder /build-staging/public/build/ ./public/build/\n")
	}

	file.WriteString(`RUN make lint-php test-php build-php-prod
RUN mkdir var
RUN chown www-data:www-data var
`)
}
