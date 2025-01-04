default: dev

format:
	@go fmt ./...

build/tailwind:
	@npm i tailwindcss@next @tailwindcss/cli@next --no-save
	@npx @tailwindcss/cli -i ./styles/styles.css -o ./.dist/css/tailwind.css --minify

dev/tailwind:
	@npm i tailwindcss@next @tailwindcss/cli@next --no-save
	@npx @tailwindcss/cli -i ./styles/styles.css -o ./.dist/css/tailwind.css --watch

dev/server:
	go run github.com/air-verse/air@latest \
	--build.cmd "go build -o tmp/main cmd/main.go" \
	--build.bin "./tmp/main" \
	--build.delay "100" \
	--build.include_ext "go,html,css,svg,yml" \
	--build.exclude_dir ".dist,node_modules" \
	--misc.clean_on_exit true

build:
	@make build/tailwind
	@go run cmd/main.go --mode=build

dev:
	@make -j2 dev/tailwind dev/server
