# tryout-grpc-go

[![standard-readme compliant](https://img.shields.io/badge/standard--readme-OK-green.svg?style=flat-square)](https://github.com/RichardLitt/standard-readme)

gRPC tryout in golang. This is based out of this [gRPC course](https://www.udemy.com/course/grpc-golang)

## Table of Contents

- [Install](#install)
- [Usage](#usage)
- [Maintainers](#maintainers)
- [Contributing](#contributing)
- [License](#license)

## Install

```
# Install Protobuf
brew install protobuf

# Install go plugins for protobuf & grpc
go install google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
```
Note: Download this source code inside the GOPATH (Havent added go mod yet)
(eg: ~/go/src/github.com/worldofprasanna/grpc-go-code)

## Usage

There are 2 gRPC servers,
1. Greet
2. Calculator

To generate the protobuf files, use the `./generator.sh` script

```
# To start the gRPC server,
go run greet/greet_server/server.go # (for Greet)
go run calculator/server/server.go # (for Calculator)

# To start the gRPC Client,
go run greet/greet_client/client.go
go run calculator/client/client.go
```

To test the gRPC server using [evans cli](https://github.com/ktr0731/evans)
```
# After installing evans cli,
evans -p 50052 -r
```

## Maintainers

[@worldofprasanna](https://github.com/worldofprasanna)

## Contributing

PRs accepted.

Small note: If editing the README, please conform to the [standard-readme](https://github.com/RichardLitt/standard-readme) specification.

## License

MIT © 2021 Prasanna V
