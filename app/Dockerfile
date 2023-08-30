FROM golang:1.21-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash make gcc musl-dev

COPY ["app/go.mod", "app/go.sum", "./"]

RUN go mod download


COPY app ./

RUN go build -o ./bin/app cmd/segment_service/main.go

FROM alpine as RUNNER

COPY --from=builder /usr/local/src/bin/app /
COPY config/local.yaml /local.yaml

CMD ["/app", "--path=/local.yaml"]