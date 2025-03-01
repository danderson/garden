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
	go run ./cmd/watch

.PHONY: live
live:
	go run ./cmd/watch

.PHONY: bin
bin:
	go build -o garden.bin -tags=osusergo,netgo

.PHONY: schema
schema:
	sqlite3 garden_dev.db ".schema"

.PHONY: sqlite
sqlite:
	sqlite3 garden_dev.db

.PHONY: deploy
deploy:
	go tool sqlc generate
	go tool templ generate
	make bin
	rsync garden.bin acrux:/fast/garden/garden.bin
	ssh root@acrux systemctl restart garden
