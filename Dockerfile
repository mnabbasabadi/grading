# Build stage
FROM golang:1.20-alpine AS builder

# Variables
ENV APP_PATH=/opt/digiexam
ENV APP_NAME=grading


WORKDIR /src

ADD . $APP_NAME

WORKDIR /src/$APP_NAME/service
RUN go build -o main ./cmd/$APP_NAME/*.go

# Build final image
FROM alpine:3.17.2

ENV APP_PATH=/opt/digiexam
ENV APP_NAME=grading

WORKDIR $APP_PATH/$APP_NAME

# Copy app binary
COPY --from=builder /src/$APP_NAME/service/main $APP_PATH/$APP_NAME/$APP_NAME


EXPOSE 8080/tcp

CMD ["/opt/digiexam/grading/grading"]
