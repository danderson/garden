.PHONY: migrate
dev:
	sqlc generate
	go run .

.PHONY: devreset
devreset:
	cp -f db.sqlite3.2 garden.db
	sqlc generate
	go run .

.PHONY: fromprod
fromprod:
	rm -f db.sqlite3.2
	fly sftp get /state/garden/db.sqlite3 db.sqlite3.2
	sqlite3 db.sqlite3.2 "PRAGMA user_version=1;"
	cp -f db.sqlite3.2 garden.db
	sqlc generate
	go run .

.PHONY: schema
schema:
	sqlite3 garden.db ".schema"
