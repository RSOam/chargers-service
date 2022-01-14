FROM golang:1.17.3-alpine3.14 as builder
RUN apk add git
RUN mkdir /chargers-service
ADD . /chargers-service
WORKDIR /chargers-service

COPY go.mod go.sum swagger.json ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

RUN mkdir /chargers-service

WORKDIR /chargers-service/

COPY --from=builder /chargers-service/swagger.json .
COPY --from=builder /chargers-service/main .

ARG DBpw_arg=default_value 
ENV DBpw=$DBpw_arg

EXPOSE 8080 50051

CMD ["./main"]