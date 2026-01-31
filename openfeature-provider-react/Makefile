.PHONY: build clean publish

build:
	npx tsc

clean:
	rm -rf dist

publish:
	npm version patch
	npm run build
	npm publish --registry http://localhost:4873