FROM golang:alpine AS builder

RUN apk add git upx

WORKDIR /code
COPY ./handler ./
COPY ./go.??? ./

RUN go build -tags container -o main . && upx main

FROM alpine
RUN apk add ca-certificates
WORKDIR /service
COPY --from=builder /code/main ./

CMD ["./main"]
