ARG BASE_IMAGE=alpine:3.18
FROM $BASE_IMAGE

RUN addgroup rpay --gid 1002 --system \
  && adduser rpay --uid 1002 --home /home/rpay --ingroup rpay --shell /bin/sh --system \
  && chown -R rpay.rpay /home/rpay \
  && apk add --no-cache --upgrade \
  curl \
  netcat-openbsd

ENV LOG_LEVEL="INFO"

USER 1002
WORKDIR /home/rpay
COPY app .

USER root
RUN chown -R rpay.rpay /home/rpay
USER 1002

EXPOSE 8085

ENTRYPOINT ["./app"]