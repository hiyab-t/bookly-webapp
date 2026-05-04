# base image
FROM golang:1.26.2

# define the working directory
WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

# build
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping ./cmd/api

# run the compiled binary when the container starts
CMD ["/docker-gs-ping"]