FROM golang:1.21 AS go

WORKDIR /build
COPY proxy proxy
COPY go.mod go.sum ./
RUN CGO_ENABLED=0 go build -v -o proxy ./proxy

FROM python:3.11-slim as py

WORKDIR /app
RUN python -m venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"
COPY ./requirements.txt ./manage.py ./
RUN pip install -Ur requirements.txt
COPY --from=go /build/proxy .
COPY garden garden
COPY boxinventory boxinventory
COPY herbarium herbarium
COPY dockerstuff/* ./
RUN ./manage.py collectstatic --noinput --clear

ENV DJANGO_SETTINGS_MODULE=garden.prod_settings

CMD ["/app/run.sh"]
