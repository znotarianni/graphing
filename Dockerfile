FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod main.go index.html ./
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /plotter .

FROM scratch
COPY --from=build /plotter /plotter
EXPOSE 8080
ENTRYPOINT ["/plotter"]

# Build:  docker build -t plotter .
# Run:    docker run --rm -p 8080:8080 -v /path/to/data:/data:ro plotter
# Open:   http://localhost:8080
