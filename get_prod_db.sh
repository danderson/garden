#!/usr/bin/env bash

rm db.sqlite3
fly sftp get /state/garden/db.sqlite3 db.sqlite3
