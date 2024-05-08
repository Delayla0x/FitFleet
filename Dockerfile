FROM golang:1.16 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

FROM node:14 as frontend
WORKDIR /app
COPY --from=builder /app/frontend .
RUN npm install
RUN npm run build

FROM debian:buster-slim
COPY --from=builder /app/main .
COPY --from=frontend /app/build /var/www/html
EXPOSE 8080
CMD ["./main"]
