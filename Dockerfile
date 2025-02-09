###############################################################################
# Build a Docker Container of the Smarter Command-line Interface.
###############################################################################
FROM golang:latest

LABEL maintainer="Lawrence McDaniel <lawrence@querium.com>"

# Create a non-root user to run the application
RUN adduser --disabled-password --gecos '' smarter_user

WORKDIR /cli

COPY go.mod go.sum ./

RUN go mod download

COPY ./cmd ./cmd
COPY main.go VERSION ./

RUN chown smarter_user:smarter_user -R .

RUN version=$(cat VERSION) && \
    CGO_ENABLED=0 GOOS=linux go build -o smarter main.go -ldflags "-X main.Version=${version}" -o ./smarter .

RUN chown smarter_user:smarter_user -R .

# setup the run-time environment
WORKDIR /
ENV PATH="/cli:${PATH}"
USER smarter_user
CMD ["/bin/sh", "-c", "while :; do sleep 10; done"]
