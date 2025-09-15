go test .
go test -v ./...
go test -v -run Test_isPrime

run groups / test suits
go test -v -run Test_app

ignorar o cache dos testes:
go test -count=1 ./...

 go test ./... -coverprofile=coverage.out -covermode=atomic \
  && go tool cover -html=coverage.out -o coverage.html \
  && (xdg-open coverage.html >/dev/null 2>&1 || gio open coverage.html || open coverage.html || start coverage.html)

go get -u github.com/ory/dockertest/v3
