FROM golang:1.26 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/cmdry . && \
    CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/cmdry-ports ./plugins/ports/cmd/cmdry-ports

FROM alpine:3.22
RUN apk add --no-cache iproute2
COPY --from=build /out/cmdry /usr/local/bin/cmdry
COPY --from=build /out/cmdry-ports /opt/cmdry/plugins/cmdry-ports
RUN mkdir -p /opt/cmdry/data
ENV CMDRY_ADDR=:8080 CMDRY_PLUGIN_DIR=/opt/cmdry/plugins CMDRY_DATA_DIR=/opt/cmdry/data
EXPOSE 8080
ENTRYPOINT ["cmdry", "serve"]
