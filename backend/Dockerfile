FROM golang:latest
WORKDIR /app

RUN go install github.com/air-verse/air@latest
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN apt-get update &&apt-get install -y tmux

# Volume mount needed for hot reloading
# Copy go.mod & go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest separately
COPY . .

EXPOSE 8081

# Use absolute path to air binary and config file
CMD ["air"]