VERSION := $(shell git describe --tags --abbrev=0)
APP := tessera

build: tailwind.css
	go build -o tessera cmd/main.go

tailwind.css:
	curl https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4 > cmd/tailwind.css

release: tailwind.css
	GOOS=linux GOARCH=amd64 \
		go build -tags release \
		-ldflags '-X github.com/joseph0x45/tessera/internal/buildinfo.Version=$(VERSION)' \
		-o $(APP)
