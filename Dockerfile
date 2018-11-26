FROM debian:stretch-slim

RUN apt-get update && \
  apt-get install -y --no-install-recommends --no-install-suggests \
  ca-certificates \
  libxml2 \
  && apt-get clean

COPY ./build/bin /bin

ENTRYPOINT ["/bin/bblfsh-json-proxy"]
CMD ["serve"]
