.PHONY: build clean deploy

build:
	npx sls build

clean:
	rm -rf ./bin

deploy:
	aws-vault exec ninedemons-admin_role -- npx sls deploy --verbose
