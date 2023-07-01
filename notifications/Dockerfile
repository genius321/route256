FROM ubuntu:22.04

RUN apt-get update && apt-get install ca-certificates -y && update-ca-certificates

ADD ./bin/app /app

CMD ["/app"]
