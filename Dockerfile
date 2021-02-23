FROM --platform=${BUILDPLATFORM} alpine:3.13.2

ARG TARGETOS
ARG TARGETARCH
ARG TARGETPLATFORM

LABEL maintainer="Gareth Evans <gareth@bryncynfelin.co.uk>"
COPY dist/uc-${TARGETOS}_${TARGETOS}_${TARGETARCH}/uc /usr/bin/uc
COPY github-actions-entrypoint.sh /usr/bin

ENTRYPOINT [ "/usr/bin/uc" ]

CMD ["--help"]
