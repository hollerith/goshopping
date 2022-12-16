# Install git as part of the build stage
FROM golang:1.15-alpine as build
RUN apk add --no-cache git

# Set the working directory
WORKDIR /app

# Copy the Go module files and download the dependencies
COPY go.mod ./
RUN go mod download

# Copy the rest of the project files
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -ldflags="-w -s" -o /app/app

# Create a minimal runtime image
FROM scratch

# Copy the compiled Go app and the .env file from the build stage
COPY --from=build /app/app /app/app
COPY --from=build /app/.env /app/.env

# Use the ARG and ENV instructions to pass the SHOPIFY_ACCESS_TOKEN value
# from the .env file to the Docker image
ARG SHOPIFY_ACCESS_TOKEN
ENV SHOPIFY_ACCESS_TOKEN=$SHOPIFY_ACCESS_TOKEN

# Expose the app's port
EXPOSE 8080

# Run the app when the container starts
CMD ["/app/app"]
