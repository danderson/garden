.PHONY: fromprod

fromprod:
	curl -o garden_prod.db http://garden/.magic/backup
	curl -o images_prod.tgz http://garden/.magic/images
	rm -f garden_dev.db*
	cp -f garden_prod.db garden_dev.db
	rm -rf images
	tar xf images_prod.tgz
	iex -S mix phx.server

reset:
	rm -f garden_dev.db*
	cp -f garden_prod.db garden_dev.db
	rm -rf images
	tar xf images_prod.tgz
	iex -S mix phx.server

dev:
	iex -S mix phx.server
