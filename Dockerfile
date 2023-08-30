FROM golang:1.21-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash make gcc musl-dev

COPY ["go.mod", "go.sum", "./"]

RUN go mod download


COPY ./ ./

RUN go build -o main cmd/segment_service/main.go

FROM alpine as RUNNER

COPY --from=builder /usr/local/src/main /
COPY config/local.yaml /local.yaml

CMD ["/main", "--path=/local.yaml"]