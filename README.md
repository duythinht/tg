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

### How to start it

* Start the all-in-one environment

```
docker-compose up
```

* Then you can test the API at port `10080` (docker map 10080 -> 8080)

* To start the API standalone

```
go run cmd/api/main.go
```

* To start the Worker standalone (be careful, the once is running on docker-compose, you might stop it before)

```
go run cmd/worker/main.go
```

* environment as config
    * DB_DSN: the DSN of postgresql
    * QUEUE_BROKERS: kafka brokers addresses
    * QUEUE_TOPIC: kafka topic name
    * QUEUE_GROUP_ID: kafka consumer group id
    * SERVER_ADDR: address for API listen, default to `8080`

### How to test is it working?

* Create a repository

```
# create a repo by ssh url

curl -XPOST -H 'Content-Type: application/json' http://localhost:10080/api/v1/repositories -d '{ "url": "git@github.com:just-a-nomad-org/repo-with-secrets.git" }'

# create a repo by https url

curl -XPOST -H 'Content-Type: application/json' http://localhost:10080/api/v1/repositories -d '{ "url": "https://github.com/just-a-nomad-org/repo-without-secrets" }'

```
output
```
{
  "repositoryId": 2
}
```

* List repositories

```
curl -XGET http://localhost:10080/api/v1/repositories
```

output

```
{
  "repositories": [
    {
      "id": 1,
      "host": "github",
      "owner": "just-a-nomad-org",
      "repository": "repo-with-secrets",
      "url": "https://github.com/just-a-nomad-org/repo-with-secrets",
      "created_at": "2022-09-09T07:00:10.138211Z",
      "updated_at": "2022-09-09T07:00:10.138211Z"
    },
    {
      "id": 2,
      "host": "github",
      "owner": "duythinht",
      "repository": "tg",
      "url": "https://github.com/duythinht/tg",
      "created_at": "2022-09-09T07:01:57.255618Z",
      "updated_at": "2022-09-09T07:01:57.255618Z"
    }
  ]
}
```

* Trigger scan for a repo
```
# eg for repositoryId = 1

curl -XPOST -H 'Content-Type: application/json' http://localhost:10080/api/v1/scans/1 -d '{}'
```

* List of scans for a repo

```
# eg for repositoryid = 1

$ curl -XGET http://localhost:10080/api/v1/scans/2

```

* Example scans outputs:

```
{
  "scans": [
    {
      "id": 1,
      "repository_id": 1,
      "status": "Failure",
      "findings": [
        {
          "type": "sast",
          "ruleId": "G143",
          "location": {
            "path": "cmd/main.go",
            "positions": {
              "begin": {
                "line": 4
              }
            }
          },
          "metadata": {
            "severity": "HIGH",
            "description": "secrets found the the source code"
          }
        }
      ],
      "queued_at": "2022-09-09T07:00:27.448445Z",
      "scanning_at": "2022-09-09T07:00:28.467063Z",
      "finished_at": "2022-09-09T07:00:29.696075Z"
    }
  ]
}
```