## Build (About 500mb image)
FROM golang:1.18-buster AS build

WORKDIR /app

# COPY go.mod .
# COPY go.sum .
# COPY *.go ./
COPY ./ ./

RUN go mod tidy

RUN go build -o /belphegorv2-build

## Deploy (Reduced to 20mb or so)
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /belphegorv2-build /belphegorv2-build

EXPOSE 8080
COPY docker-compose.yml .en[v] ./
COPY ./sounds_assets/ ./sounds_assets/

USER nonroot:nonroot

ENTRYPOINT ["/belphegorv2-build"]
