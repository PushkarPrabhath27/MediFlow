FROM golang:1.25-alpine

WORKDIR /app

# Install air for hot reloading
RUN go install github.com/air-verse/air@latest

# Install dependencies separately for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Expose the port
EXPOSE 8080

# Command to run air
CMD ["air", "-c", ".air.toml"]
