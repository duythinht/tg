# Github Repository Scan

### Migration

* Create migration

```bash
migrate -dir ./migrations -ext sql <name of migration>
```

* Run the migration
```
migrate -path ./migrations -database postgres://postgres:x@localhost:5432/postgres?sslmode=disable up
```