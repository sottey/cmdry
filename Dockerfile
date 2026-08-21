FROM golang:1.26 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN version=$(tr -d '\r\n' < VERSION) && \
    linker_flags="-s -w -X github.com/sottey/cmdry/internal/buildinfo.Version=$version -X github.com/sottey/cmdry/plugin-sdk/go.BuildVersion=$version" && \
    CGO_ENABLED=0 go build -trimpath -ldflags="$linker_flags" -o /out/cmdry . && \
    CGO_ENABLED=0 go build -trimpath -ldflags="$linker_flags" -o /out/cmdry-ports ./plugins/ports/cmd/cmdry-ports && \
    CGO_ENABLED=0 go build -trimpath -ldflags="$linker_flags" -o /out/cmdry-journal ./plugins/journal/cmd/cmdry-journal

FROM alpine:3.22
RUN apk add --no-cache iproute2
COPY --from=build /out/cmdry /usr/local/bin/cmdry
COPY --from=build /out/cmdry-ports /opt/cmdry/plugins/cmdry-ports
COPY --from=build /out/cmdry-journal /opt/cmdry/plugins/cmdry-journal
RUN mkdir -p /opt/cmdry/data
ENV CMDRY_ADDR=:8080 CMDRY_PLUGIN_DIR=/opt/cmdry/plugins CMDRY_DATA_DIR=/opt/cmdry/data
EXPOSE 8080
ENTRYPOINT ["cmdry", "serve"]
