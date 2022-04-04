all: build

build:
	go build -v

install:
	go install -v

spell:
	npx spellchecker-cli *.txt *.md LICENSE -d .dictionary
