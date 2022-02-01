FROM golang:1.17-alpine AS base

ENV GO111MODULE=on \
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o tildly .

WORKDIR /dist

RUN cp /build/tildly .

FROM scratch

COPY --from=base /dist/tildly /
COPY ./handler/templates/ /handler/templates

EXPOSE 8080

CMD [ "/tildly" ]
