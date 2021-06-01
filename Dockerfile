#
# Build frontend
#
FROM node:15-alpine as frontend-builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY ./web ./web
COPY ./*.config.js .
RUN npm run build

#
# Build backend
#
FROM golang:1.16-alpine as backend-builder

ARG TARGETARCH
ARG TARGETOS

ENV GO111MODULE on
ENV GOPATH /

RUN mkdir /app && cd /app
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download && go mod verify

COPY . .
COPY --from=frontend-builder /app/web/dist /app/web/dist
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="-w -s" -o /app/nuts-demo-ehr

#
# Runtime
#
FROM alpine:3.13
RUN mkdir /app && cd /app
WORKDIR /app
COPY --from=backend-builder /app/nuts-demo-ehr .
HEALTHCHECK --start-period=5s --timeout=5s --interval=5s \
    CMD wget --no-verbose --tries=1 --spider http://localhost:1304/ || exit 1
EXPOSE 1304
ENTRYPOINT ["/app/nuts-demo-ehr"]