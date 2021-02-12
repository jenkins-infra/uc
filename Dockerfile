FROM --platform=${BUILDPLATFORM} alpine:3.13.1

ARG TARGETOS
ARG TARGETARCH
ARG TARGETPLATFORM

LABEL maintainer="Gareth Evans <gareth@bryncynfelin.co.uk>"
COPY dist/uc-${TARGETOS}_${TARGETOS}_${TARGETARCH}/uc /usr/bin/uc

ENTRYPOINT [ "/usr/bin/uc" ]

CMD ["--help"]
