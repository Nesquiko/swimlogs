FROM golang:1.20-alpine as builder

WORKDIR /swimlogs

COPY go.mod go.sum ./

# install dependencies
RUN apk add make
RUN go mod download
RUN go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.13.0

COPY . .

RUN make generate

RUN CGO_ENABLED=0 GOOS=linux go build -o /swimlogs-app

FROM alpine:3.18

WORKDIR /

COPY --from=builder /swimlogs/migrations /migrations
COPY --from=builder /swimlogs-app /swimlogs-app

EXPOSE 42069

USER nonroot:nonroot

CMD ["/swimlogs-app"]
