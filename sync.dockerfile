FROM golang:1.16-buster AS build

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY cmd/cron/*.go ./
COPY internal/ internal/
COPY api/ api/

RUN go build -o /sync-db

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /sync-db /sync-db

USER nonroot:nonroot

ENTRYPOINT ["/sync-db"]