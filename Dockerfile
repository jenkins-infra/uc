FROM alpine:3.12.3

COPY ./build/linux/uc /usr/local/bin
COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
