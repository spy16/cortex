## Build binary
FROM golang:1.18-buster AS build
WORKDIR /code
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./ ./
RUN make

## Release Image
FROM gcr.io/distroless/base-debian10
WORKDIR /app
COPY --from=build /code/bin/cortex /app/
COPY --from=build /code/cortex.yaml /app/
EXPOSE 8080
ENV PORT=8080
USER nonroot:nonroot
CMD ["/app/cortex", "serve"]
