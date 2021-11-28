ARG GO_VERSION=1.16
FROM golang:${GO_VERSION}-buster AS dev

RUN apt install -y ca-certificates git

ENV APP_NAME="backend" \
    APP_PATH="/var/app"

ENV APP_BUILD_NAME="${APP_NAME}"

COPY . ${APP_PATH}
WORKDIR ${APP_PATH}

ENV GO111MODULE="on" \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOFLAGS="-mod=vendor"

ENTRYPOINT ["sh"]

FROM dev as build

RUN (([ ! -d "${APP_PATH}/vendor" ] && go mod download && go mod vendor) || true)
RUN GIT_COMMIT=$(git rev-list -1 HEAD) && \
    BUILD=$(date +%FT%T%z) && \
    go build -ldflags="-s -w -X 'main.Version=${GIT_COMMIT}' -X main.Build=${BUILD}" -mod vendor -o ${APP_BUILD_NAME} .
RUN chmod +x ${APP_BUILD_NAME}

FROM debian:10-slim AS prod
RUN apt update -y && apt install -y curl iproute2

ENV APP_BUILD_PATH="/var/app" \
    APP_BUILD_NAME="backend"

WORKDIR ${APP_BUILD_PATH}

COPY --from=build ${APP_BUILD_PATH}/${APP_BUILD_NAME} ${APP_BUILD_PATH}/

ENTRYPOINT ["/var/app/backend"]
CMD ""