.PHONY: migrate
dev:
	./manage.py makemigrations
	./manage.py migrate
	./manage.py runserver 0.0.0.0:8000

.PHONY: devreset
devreset:
	cp -f db.sqlite3.2 db.sqlite3
	./manage.py makemigrations
	./manage.py migrate
	./manage.py runserver 0.0.0.0:8000

.PHONY: fromprod
fromprod:
	rm -f db.sqlite3
	fly sftp get /state/garden/db.sqlite3 db.sqlite3.2
	cp -f db.sqlite3.2 db.sqlite3
	./manage.py makemigrations
	./manage.py migrate
	./manage.py runserver 0.0.0.0:8000
