package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/worldofprasanna/grpc-go-code/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

type server struct{
	greetpb.UnimplementedGreetServiceServer
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function is invoked with %v", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := greetpb.GreetResponse{
		Result: result,
	}
	return &res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function is invoked with %v", req)
	firstName := req.GetGreeting().FirstName

	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number: " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {

	fmt.Printf("LongGreet function is invoked with stream\n")
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		} else if err != nil {
			log.Fatalf("Error while reading client stream %v", err)
		} else {
			firstName := req.Greeting.FirstName
			result += "Hello " + firstName + "! "
		}
	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("GreetEveryone function is invoked with stream\n")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		} else if err != nil {
			log.Fatalf("Error while readign client stream %v", err)
			return err
		}
		firstName := req.GetGreeting().FirstName
		result := "Hello " + firstName
		err = stream.Send(&greetpb.GreetEveryoneResponse{
			Response: result,
		})
		if err != nil {
			log.Fatalf("Error while sending message to client %v", err)
			return err
		}
	}
}

func (*server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetRequestWithDeadline) (*greetpb.GreetResponseWithDeadline, error) {
	fmt.Printf("Greet with deadline function is invoked with %v\n", req)
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			fmt.Println("Client cancelled the request")
			return nil, status.Error(codes.Canceled, "Client cancelled the request")
		}
		time.Sleep(1 * time.Second)
	}
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := greetpb.GreetResponseWithDeadline{
		Response: result,
	}
	return &res, nil
}

func main()  {
	fmt.Println("Hello World")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen on 50051 %w", err)
	}

	tls := true
	opts := []grpc.ServerOption{}

	if tls {
		certFile := "ssl/server.crt"
		keyFile := "ssl/server.pem"

		creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
		if sslErr != nil {
			log.Fatalf("Failed to load certificates %v", sslErr)
			return
		}
		opts = append(opts, grpc.Creds(creds))
	}

	s := grpc.NewServer(opts...)

	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}