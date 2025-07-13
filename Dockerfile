FROM golang:1.24.5

WORKDIR /app

RUN go install github.com/air-verse/air@latest \
    && go install golang.org/x/tools/cmd/goimports@latest \
    && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
       | sh -s -- -b /go/bin

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

CMD ["air"]
