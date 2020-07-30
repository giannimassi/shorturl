# shorturl

Another url shortener service.

## Description

A web server providing a url shortener service which serves redirected urls on `/` and an API for adding, deleting and getting information on short urls. The API documentation is additionaly provided via an additional endpoint

## Usage

The project uses a _Makefile_ with the following commands:

- `make gen`: calls go generate, building the swagger documentation with [`swag`](https://github.com/swaggo/swag)
- `make test`: runs all tests in the source code
- `make lint`: runs the linter on the source code
- `make run`: runs the web server on port

To add a new short url:

```bash
curl --header "Content-Type: application/json" --request PUT --data '{"Key":"a", "URL":"http://example.org/a"}' http://localhost:8080/api -v
```

To see the full API documentation start the server and go to [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html).

### Requirements for building and generating documentation

Install [`swag`](https://github.com/swaggo/swag) with the following command:

```bash
go get -u github.com/swaggo/swag/cmd/swag
```

-------
_TODO_

- Refactor handlers to use `gin.Handler` signature
- Use gin for testing mux
- persistent storage
- docker-compose-based end-to-end/integration testing
