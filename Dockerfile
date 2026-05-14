FROM golang:1.26

WORKDIR /src

COPY . .
RUN go mod download
RUN apt-get update && apt-get install -y postgresql-client

RUN go build -o /src/main ./app/cmd

EXPOSE 8080

CMD ["/src/main"]
