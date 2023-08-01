FROM golang as builder

WORKDIR /app
# COPY go.mod /app
# COPY go.sum /app

# RUN go mod download

# COPY . /app
COPY . .

RUN go build -o sitemanagerservice /app/main.go

# FROM alpine:latest

# WORKDIR /app

# COPY --from=builder /app/ /app/

EXPOSE 8082

CMD ./sitemanagerservice