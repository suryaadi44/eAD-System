FROM golang:1.19-alpine3.16 as builder
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o /bin/main -v ./cmd/main

FROM surnet/alpine-wkhtmltopdf:3.16.0-0.12.6-small as wkhtmltopdf

FROM alpine:3.16

# Install dependencies for wkhtmltopdf
RUN apk add --no-cache \
  libstdc++ \
  libx11 \
  libxrender \
  libxext \
  libssl1.1 \
  ca-certificates \
  fontconfig \
  freetype \
  ttf-dejavu \
  ttf-droid \
  ttf-freefont \
  ttf-liberation

# Copy wkhtmltopdf files from docker-wkhtmltopdf image
COPY --from=wkhtmltopdf /bin/wkhtmltopdf /bin/wkhtmltopdf

# Copy the binary from the builder image
COPY --from=builder /bin/main .
COPY --from=builder /app/template/signature ./template/signature

CMD [ "./main" ]
