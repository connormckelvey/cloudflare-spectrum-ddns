FROM golang:1.17-alpine as build

WORKDIR /build

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /update-spectrum-ip ./cmd/... 


FROM alpine:3.15.0

WORKDIR /app

COPY --from=build /update-spectrum-ip .
COPY docker/docker-entrypoint.sh .

ENTRYPOINT [ "/app/docker-entrypoint.sh" ]
