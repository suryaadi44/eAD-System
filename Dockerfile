FROM golang:1.19.2-alpine3.16 as builder
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o /bin/main -v ./cmd/main

FROM alpine:3.16
COPY --from=builder /bin/main .
COPY --from=builder /app/template/signature ./template/signature
CMD [ "./main" ]
