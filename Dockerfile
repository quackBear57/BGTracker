# syntax=docker/dockerfile:1

FROM golang:1.27-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/bgtracker ./cmd/web

FROM alpine:3.22

RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -H -u 10001 bgtracker

WORKDIR /app

COPY --from=build /out/bgtracker /usr/local/bin/bgtracker
COPY assets ./assets

USER bgtracker

EXPOSE 8080

ENTRYPOINT ["bgtracker"]
CMD ["-addr", ":8080"]
