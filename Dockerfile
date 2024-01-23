FROM alpine:latest

WORKDIR /app/

COPY ./db-table-metrics /usr/local/bin/db-table-metrics

USER 30463

ENTRYPOINT [ "/usr/local/bin/db-table-metrics" ]