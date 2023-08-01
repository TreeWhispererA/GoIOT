FROM golang as builder

WORKDIR /app
# COPY go.mod /app
# COPY go.sum /app

# RUN go mod download

# COPY . /app
COPY . .

RUN go build -o devicemanagerservice /app/main.go

# FROM alpine:latest

# WORKDIR /app

# COPY --from=builder /app/ /app/

EXPOSE 8083

CMD ./devicemanagerservice