# HelloApp

HelloApp is a simple HTTP server that is meant to simply say "hello" when you
make a `GET` request.

## Building

You can run the following from the project root to build HelloApp:

    $ mkdir -p dist
    $ go build -o dist/helloapp ./cmd/helloapp

## Usage

    Usage of helloapp:
      -port uint
            Port for service to listen on (default 8080)


Simply run the `helloapp` executable after building and it will listen on the
default port:

    $ ./dist/helloapp
    2019/02/14 20:44:50 helloapp: starting
    2019/02/14 20:44:50 helloapp: my hostname is: myhost.local
    2019/02/14 20:44:50 helloapp: listening on port: 8080

Press `Ctrl+C` to stop the server.

## Testing

This project makes use of [Ginkgo](https://onsi.github.io/ginkgo/) for testing.

Ideally, use the `ginkgo` command-line tool for testing:

    $ ginkgo -v ./cmd/helloapp

Alternately, you can use `go test` like so:

    $ go test ./cmd/helloapp
