.PHONY: dev build test lint clean docker-build

dev:
	pnpm dev

build:
	pnpm build

preview:
	pnpm preview

test:
	pnpm test

lint:
	pnpm lint

typecheck:
	pnpx tsc --noEmit

clean:
	rm -rf dist/ node_modules/

docker-build:
	docker build -t industrial-iot/web-portal .
