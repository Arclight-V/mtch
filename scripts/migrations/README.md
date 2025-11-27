# database migrations (`migrate`)

for database migrations, it is used [migrate](https://github.com/golang-migrate/migrate)
[PostgreSQL tutorial for beginners](https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md)

## Configure database
```bash
export POSTGRESQL_URL='postgres://{user}:{password}@localhost:5432/{db_name}?sslmode=disable'
```

## Run migrations
```bash
migrate -database ${POSTGRESQL_URL} -path db/migrations up
```


## migrate db commands (`makefile`)

To run migration scripts, use the rules in the makefile.

## example
```bash
make -f scripts/migrations/Makefile migrate-up
```

