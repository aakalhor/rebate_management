# Use an official Golang image as the base
FROM golang:1.22

# Set the working directory
WORKDIR /app
RUN rm -rf docs
RUN go install github.com/swaggo/swag/cmd/swag@latest
#RUN go install github.com/aws/aws-sdk-go-v2/config@latest
#RUN go install github.com/aws/aws-sdk-go-v2/service/dynamodb@latest


# Copy the Go modules and source code
COPY go.mod go.sum ./
RUN go mod download
COPY . .



# Create the directory for the SQLite database file

RUN rm -rf ./docs && swag init --parseDependency --parseInternal
RUN mkdir -p /data

# Command to run the application directly
CMD ["go", "run", "main.go"]
