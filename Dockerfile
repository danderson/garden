FROM golang:1.20 AS go

WORKDIR /build
COPY ./proxy .
RUN CGO_ENABLED=0 go build -v -o proxy .

FROM python:3-slim as py

WORKDIR /app
RUN python -m venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"
COPY ./requirements.txt ./manage.py ./
RUN pip install -Ur requirements.txt
COPY --from=go /build/proxy .
COPY garden garden
COPY boxinventory boxinventory
COPY dockerstuff/* ./
RUN ./manage.py collectstatic --noinput --clear

ENV DJANGO_SETTINGS_MODULE=garden.prod_settings

CMD ["/app/run.sh"]
