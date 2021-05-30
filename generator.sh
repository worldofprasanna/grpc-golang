#!/bin/bash

protoc -I=greet/greetpb --go_out=greet --go-grpc_out=greet greet.proto
protoc -I=calculator/calculatorpb --go_out=calculator --go-grpc_out=calculator calculator.proto