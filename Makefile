.PHONY: fromprod
fromprod:
	curl -o garden_prod.db http://garden/.magic/db
	curl -o images_prod.tgz http://garden/.magic/images
	make reset

.PHONY: reset
reset:
	rm -f garden_dev.db*
	cp -f garden_prod.db garden_dev.db
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

.PHONY: bin
bin:
	go build -o garden.bin -tags=osusergo,netgo,sqlite_stat4,sqlite_omit_load_extension,sqlite_fts5,sqlite_json,sqlite_math_functions -ldflags=-extldflags=-static

.PHONY: schema
schema:
	sqlite3 garden_dev.db ".schema"

.PHONY: sqlite
sqlite:
	sqlite3 garden_dev.db

.PHONY: deploy
deploy:
	sqlc generate
	templ generate
	tailwindcss -i style.css -o static/app.css
	make bin
	rsync garden.bin acrux:/fast/garden/garden.bin
	ssh root@acrux systemctl restart garden
