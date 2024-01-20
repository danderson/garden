.PHONY: fromprod
fromprod:
	curl -o garden_prod.db http://garden/.magic/backup
	curl -o images_prod.tgz http://garden/.magic/images
	make reset

.PHONY: reset
reset:
	rm -f garden_dev.db*
	cp -f garden_prod.db garden_dev.db
	sqlite3 garden_dev.db "pragma user_version=1"
	rm -rf images
	tar xf images_prod.tgz
	make live

.PHONY: dev
dev:
	sqlc generate
	templ generate
	tailwindcss -i style.css -o static/app.css
	go run .

.PHONY: live
live:
	go run ./cmd/watch

.PHONY: schema
schema:
	sqlite3 garden_dev.db ".schema"

.PHONY: sqlite
sqlite:
	sqlite3 garden_dev.db
