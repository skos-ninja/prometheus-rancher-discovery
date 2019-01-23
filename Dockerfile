FROM golang:1.11.3-alpine as builder
ENV CGO_ENABLED=0
RUN mkdir -p /usr/local/app
WORKDIR /usr/local/app
COPY ./go.mod .
# COPY ./go.sum .
RUN go mod download
COPY . /usr/local/app
RUN go build .

FROM alpine
ENV RANCHER_HOST=
ENV DEFAULT_FILE=/usr/local/output
RUN mkdir -p /usr/local/app
RUN mkdir -p /usr/local/output
WORKDIR /usr/local/app
ENTRYPOINT ["./discovery"]
COPY --from=builder /usr/local/app/discovery .
