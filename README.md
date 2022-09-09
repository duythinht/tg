# Github Repository Scan

## Prerequisite
* Go 1.19: project is building with go1.19, but you also can use go1.15++ as well
* Postgres 14: Included on docker-compose
* Kafka 2.x: Included on docker-compose
* An Editor: to write/read this codebase

## Tech Stack

* echo: Just a simple http router/micro framework to handle http API
* pgx: A pure PostgreSQL driver and toolkit, which can handle postgres specific features as well
* kafka-go: Kafka library in go, to use as queued for lazy scanning and fast response on the API
* golang-migrate: CLI to create & run the database migration 
* Docker: To run develop environment
* PostgreSQL: RDBMS
* Kafka/Redpanda: Kafka is used as message broker, for develop evironment I use redpanda as an alternation for low footprint resource

## Describe the task

* Project structure

```
.
├── cmd
│   ├── api # main package for API entrypoint
│   └── worker # main package for worker entrypoint
├── api # API implement
├── worker # worker implement
├── docker-compose.yaml # set up for development env
├── go.mod
├── go.sum
├── lib # some people use `pkg` but `lib` is more relevant
│   ├── repository # package for abstraction git host
│   ├── rule # rule define & it own implement (current only secrets rule was here)
│   ├── store # store (db/queue) abstraction, useful for inject testing later
│   └── zipfs # zipfs
├── migrations # migrations scripts
├── model # database model
├── test.http # rest client foot print
└── config # all of configuration
```

* zipfs & git archive
    * I use git archive instead of clone the repo to local because:
        * the git tree is very big even you ref source tree is just 1MB, you might to fetch 1GB .git tree for large repo
        * mostly all git hosting provide the ref as archive (github, gitlab, bitbucket...)
        * We don't need local volume storage to stored the source tree
    * zipfs is abstraction on top of zipfile, we can access/walk as `io/fs` interface, thanks to go1.16++
* rule was written as an interface signature, which is easy to extensible more, currently only `secrets` was implemented
* secrets rule is implmented by lexical analysis scanner, which have low cost abstraction, the implement is determine to scan only `.go` source code (only checkGoSource method was implemented).

* Kafka producer & consumer: We need to fast response to the user whenever an scan was triggered, so kafka was used for this case.

### How to

* Run the migration
```
migrate -path ./migrations -database postgres://postgres:x@localhost:5432/postgres?sslmode=disable up
```