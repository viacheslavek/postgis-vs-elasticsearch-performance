FROM postgres:latest

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        build-essential \
        ca-certificates \
        postgresql-server-dev-all \
        postgresql-common \
        libproj-dev \
        libjson-c-dev \
        libgdal-dev \
        gdal-bin \
        postgis \
    && rm -rf /var/lib/apt/lists/*

ENV POSTGRES_USER postgresgis
ENV POSTGRES_PASSWORD passwordforgis

COPY initdb /docker-entrypoint-initdb.d
