build: reuse.go options.go stats.go
	go build -ldflags "-X main.version=0.0.`git log -1 --format=%ct`" reuse

install: reuse.go options.go stats.go
	go install -ldflags "-X main.version=0.0.`git log -1 --format=%ct`" reuse
clean:
	go clean reuse
