.PHONY: sync-content dev test build clean

# Copy frontend content into embedded/ for go:embed.
# Required before `go build` / `go run` / `go test` because go:embed needs
# the referenced paths to exist at compile time.
sync-content:
	@rm -rf internal/chat/embedded
	@mkdir -p internal/chat/embedded/projects
	@cp ../frontend/src/content/resume/profile.md internal/chat/embedded/resume.md
	@cp ../frontend/src/content/projects/*.md internal/chat/embedded/projects/
	@echo "Content synced into internal/chat/embedded/"

dev: sync-content
	go run ./cmd/server

test: sync-content
	go test ./... -v

build: sync-content
	go build -o bin/server ./cmd/server

clean:
	@rm -rf internal/chat/embedded bin
