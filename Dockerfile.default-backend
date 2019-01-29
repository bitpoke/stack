# Build the default-backend binary
FROM golang:1.11.5 as builder

# Copy in the go src
WORKDIR /go/src/github.com/presslabs/stack
COPY default-backend/ default-backend/
COPY vendor/ vendor/

# Pack templates
RUN cd default-backend && go run ../vendor/github.com/gobuffalo/packr/v2/packr2

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /default-backend github.com/presslabs/stack/default-backend

# Copy the dashboard binary into a thin image
FROM scratch
COPY --from=builder /go/src/github.com/presslabs/stack/default-backend/rootfs /
COPY --from=builder /default-backend /
ENTRYPOINT ["/default-backend"]
