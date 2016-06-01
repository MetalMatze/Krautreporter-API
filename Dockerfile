FROM alpine:latest

WORKDIR /app

RUN apk update
RUN apk add bash ca-certificates

COPY database/ database/
ADD api .

ENTRYPOINT ["/app/api"]
