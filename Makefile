.PHONY: build clean deploy

build:
	npx sls build

clean:
	rm -rf ./bin

deploy:
	aws-vault exec ninedemons-serverless -- npx sls deploy --verbose
