FROM golang:1.18-alpine
LABEL maintainer="andrey.malkevich@gmail.com"
WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o ./bin/app main.go

EXPOSE 8080

CMD ["./bin/app"]


