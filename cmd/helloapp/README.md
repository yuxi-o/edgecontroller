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

## Copyright and License

Copyright 2019 Smart-Edge.com, Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
