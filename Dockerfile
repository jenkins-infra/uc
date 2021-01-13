FROM alpine:3.12.3

COPY ./build/linux/jcasc-validator /usr/local/bin
RUN jcasc-validator --version
