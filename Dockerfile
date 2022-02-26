FROM golang:1-alpine

WORKDIR /app
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.9.0/wait /wait
RUN chmod +x /wait
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o /starwars
CMD /wait && /starwars