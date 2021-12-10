.PHONY: build clean deploy

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/bootstrap hello/main.go
	zip -j bin/hello.zip bin/bootstrap

clean:
	rm -rf ./bin

deploy:
	aws-vault exec ninedemons-admin_role -- npx sls deploy --verbose

all: clean build deploy
