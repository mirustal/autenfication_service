FROM golang:latest as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main .

FROM alpine:latest 
WORKDIR /app/


COPY config.yml  .
ENV CONFIG_PATH="./config.yml"
ENV SECRET_KEY="Medods_Task1"

COPY --from=builder /app/main .

RUN chmod +x main


EXPOSE 8081


CMD ["./main"]