# Pulling the official postgres base image
FROM postgres:16.3-alpine3.20

# Argument pathtoservice
ARG PATHTOSERVICE=./schema

# Naming default database, user, and password
ENV POSTGRES_DB postgres
ENV POSTGRES_USER pombnk
ENV POSTGRES_PASSWORD postgres

# Executing SQL file on startup
COPY $PATHTOSERVICE/init.sql /docker-entrypoint-initdb.d/