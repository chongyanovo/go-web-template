.PHONY: docker
docker:
	@rm go-web-app || true
	@GOOS=linux GOARCH=arm go build -o go-web-app cmd/main.go
    @docker rmi -f chongyan/go-web-app:v1.0
	@docker build -t chongyan/go-web-app:v1.0 .