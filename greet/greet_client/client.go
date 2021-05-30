package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/worldofprasanna/grpc-go-code/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hello I am client")

	tls := true
	opts := grpc.WithInsecure()

	if tls {
		certFile := "ssl/ca.crt"
		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
		if sslErr != nil {
			log.Fatalf("Error while loading CA Trust certificate: %v", sslErr)
			return
		}
		opts = grpc.WithTransportCredentials(creds)
	}

	conn, err := grpc.Dial("localhost:50051", opts)

	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer conn.Close()
	c := greetpb.NewGreetServiceClient(conn)
	doUnary(c)
	// doServerStreaming(c)
	// doClientStreaming(c)
	// doBiDiStreaming(c)
	// doUnaryWithDeadline(c, 5 * time.Second)
	// doUnaryWithDeadline(c, 1 * time.Second)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do unary RPC")
	request := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Prasanna",
			LastName: "V",
		},
	}
	res, err := c.Greet(context.Background(), request)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}

	log.Printf("Response from Greet: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Startiing to do server streaming")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Prasanna",
			LastName: "V",
		},
	}
	stream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes RPC %v", req)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream %v", err)
		}
		log.Printf("response from greet many times %v", msg)
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do client streaming")

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreet: %v", err)
	}

	requests := []*greetpb.LongGreetRequest {
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Prasanna",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Prasanna 1",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Prasanna 2",
			},
		},
	}

	for _, req := range requests {
		fmt.Printf("Sending Req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receving response from server")
	}

	fmt.Printf("Response from server %v\n", res)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting a Bi Di Streaming")

	requests := []*greetpb.GreetEveryoneRequest {
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Prasanna",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Prasanna 1",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Prasanna 2",
			},
		},
	}

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error while connecting with server %v", err)
		return
	}

	waitc := make(chan struct{})

	go func() {
		for _, req := range requests {
			fmt.Printf("Sending the message %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("Closing the client streaming")
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving %v", err)
				break
			}
			fmt.Printf("Received message: %v\n", resp.GetResponse())
		}
		close(waitc)
	}()

	<- waitc
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("Starting to do unary with deadline RPC")
	request := &greetpb.GreetRequestWithDeadline{
		Greeting: &greetpb.Greeting{
			FirstName: "Prasanna",
			LastName: "V",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, request)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Deadline exceeded")
			} else {
				fmt.Printf("Some other error %v\n", err)
			}
		} else {
			log.Fatalf("error while calling GreetWithDeadline %v", err)
		}
		log.Fatalf("error while calling Greet with Deadline RPC: %v", err)
		return
	}

	log.Printf("Response from Greet: %v", res.GetResponse())
}