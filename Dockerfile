FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./

RUN go mod download

RUN go mod tidy

COPY . .

RUN go build 

EXPOSE 8080

CMD [ "./tildly" ]
