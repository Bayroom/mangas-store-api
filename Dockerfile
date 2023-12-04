# Bring latest official go image
FROM golang:latest

# Define working directory
WORKDIR /app

# Copy go special files to WD
COPY mangas-store/go.mod .
COPY mangas-store/go.sum .

# Download deps
RUN go mod download

# Copy sources to WD
COPY mangas-store .

# Build app
RUN go build -o main .

# Expose app port
EXPOSE 8080

# Launch
CMD [ "./main" ]