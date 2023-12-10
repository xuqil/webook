.PHONY: docker
docker:
	@rm webook || true
	@go mod tidy
	@GOOS=linux GOARCH=arm go build -o webook .
	@docker image rmi -f xuqil/webook:v0.0.1
	@docker build -t xuqil/webook:v0.0.1 .

mock:
	@mockgen -source=internal/service/user.go -package=svcmocks -destination=internal/service/mocks/user.mock.go
	@mockgen -source=internal/service/code.go -package=svcmocks -destination=internal/service/mocks/code.mock.go
	@mockgen -source=internal/repository/code.go -package=repomocks -destination=internal/repository/mocks/code.mock.go
	@mockgen -source=internal/repository/user.go -package=repomocks -destination=internal/repository/mocks/user.mock.go
	@mockgen -source=internal/repository/dao/user.go -package=daomocks -destination=internal/repository/dao/mocks/user.mock.go
	@mockgen -source=internal/repository/cache/user.go -package=cachemocks -destination=internal/repository/cache/mocks/user.mock.go
	@mockgen -package=redismocks -destination=internal/repository/cache/redismocks/cmdable.mock.go github.com/redis/go-redis/v9 Cmdable
	@go mod tidy