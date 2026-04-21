FROM golang:1.20-alpine AS build
WORKDIR /src
COPY . .
RUN apk add --no-cache git
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/bootstrap ./cmd/api

FROM scratch
COPY --from=build /out/bootstrap /bootstrap
ENTRYPOINT ["/bootstrap"]
