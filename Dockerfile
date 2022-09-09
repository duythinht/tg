FROM golang:1.18.2 as build
WORKDIR /opt/tg/src
COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./
RUN go build -o /opt/tg/bin/api ./cmd/api/main.go
RUN go build -o /opt/tg/bin/worker ./cmd/worker/main.go

FROM migrate/migrate:4 as migrate-tool

FROM ubuntu:20.04
RUN apt-get update && apt-get install -y ca-certificates 
WORKDIR /opt/tg
COPY --from=build /opt/tg/bin /opt/tg/bin
COPY --from=migrate-tool /usr/local/bin/migrate /opt/tg/bin
COPY ./config/config.yaml /opt/tg/config.yaml
COPY ./entrypoint.sh /opt/tg/
COPY ./migrations /opt/tg/migrations
EXPOSE 8080
ENTRYPOINT [ "/opt/tg/entrypoint.sh" ]
CMD ["/opt/tg/bin/api"]