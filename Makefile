.PHONY: migrate
dev:
	./manage.py makemigrations
	./manage.py migrate
	./manage.py runserver

.PHONY: fromprod
fromprod:
	rm -f db.sqlite3
	fly sftp get /state/garden/db.sqlite3 db.sqlite3
	./manage.py makemigrations
	./manage.py migrate
	./manage.py runserver