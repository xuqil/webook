.PHONY: docker
docker:
	@rm webook || true
	@go mod tidy
	@GOOS=linux GOARCH=arm go build -o webook .
	@docker image rmi -f xuqil/webook:v0.0.1
	@docker build -t xuqil/webook:v0.0.1 .
