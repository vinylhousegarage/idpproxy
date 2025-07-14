FROM golang:1.24.5

WORKDIR /app

RUN go install github.com/air-verse/air@v1.62.0 \
    && go install golang.org/x/tools/cmd/goimports@v0.35.0 \
    && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
       | sh -s -- -b /go/bin v1.64.8

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

CMD ["air"]
