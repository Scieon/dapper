FROM postgres:latest

# Setup db env var
ENV POSTGRES_DB=dapper
ENV POSTGRES_USER=labs
ENV POSTGRES_PASSWORD=dapper

# Init DB
COPY ./init.sql /docker-entrypoint-initdb.d/
