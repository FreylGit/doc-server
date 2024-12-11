FROM alpine:3.13

RUN apk update && \
    apk upgrade && \
    apk add bash && \
    rm -rf /var/cache/apk/*

ADD https://github.com/pressly/goose/releases/download/v3.23.0/goose_linux_x86_64 /bin/goose
RUN chmod +x /bin/goose

WORKDIR /root

# Добавление необходимых файлов
ADD migration.sh .
ADD local.env .

RUN chmod +x migration.sh

ENTRYPOINT ["bash", "migration.sh"]