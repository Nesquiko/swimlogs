FROM golang:1.21-alpine as oapi-codegen-download

RUN go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.1

FROM node:21-alpine as api-gen

WORKDIR /api

COPY --from=oapi-codegen-download /go/bin/oapi-codegen /

ENV PATH="${PATH}:/"

COPY ./api /api

WORKDIR /api

RUN npm run generate-go

FROM golang:1.21-alpine as builder

COPY --from=api-gen /api /api

WORKDIR /api

RUN go mod download

WORKDIR /swimlogs

COPY ./backend/go.mod ./backend/go.sum ./

RUN go mod download

COPY ./backend ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /swimlogs-app

FROM alpine:3.19

WORKDIR /

COPY --from=builder /swimlogs/migrations /migrations
COPY --from=builder /swimlogs-app /swimlogs-app

EXPOSE 42069

HEALTHCHECK --interval=5s --timeout=3s --retries=3 \
	CMD wget --tries=1 --spider http://localhost:42069/monitoring/heartbeat || exit 1

CMD ["/swimlogs-app"]
