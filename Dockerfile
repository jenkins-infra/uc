FROM --platform=${BUILDPLATFORM} curlimages/curl:7.82.0 AS build-stage0

ARG TARGETOS
ARG TARGETARCH
ARG TARGETPLATFORM
ARG UC_VERSION

RUN curl -L -o /tmp/uc-${TARGETOS}-${TARGETARCH}.tar.gz https://github.com/jenkins-infra/uc/releases/download/${UC_VERSION}/uc-${TARGETOS}-${TARGETARCH}.tar.gz && \
      tar -xvzf /tmp/uc-${TARGETOS}-${TARGETARCH}.tar.gz -C /tmp && \
      chmod a+x /tmp/uc

FROM --platform=${BUILDPLATFORM} alpine:3.15.4
LABEL maintainer="Gareth Evans <gareth@bryncynfelin.co.uk>"

COPY --from=build-stage0 /tmp/uc /usr/bin/uc
COPY github-actions-entrypoint.sh /usr/bin

ENTRYPOINT [ "/usr/bin/uc" ]
CMD ["--help"]
