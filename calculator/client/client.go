package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/worldofprasanna/grpc-go-code/calculator/calculator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hello from Client")

	conn, err := grpc.Dial("0.0.0.0:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error connecting to server %v", err)
	}

	defer conn.Close()

	service := calculator.NewCalculatorServiceClient(conn)
	// doUnary(service)
	// doServerStreaming(service)
	// doClientStreaming(service)
	// doBidirectionalStreaming(service)
	doErrorUnary(service)
}

func doUnary(client calculator.CalculatorServiceClient) {
	request := &calculator.MathRequest{
		FirstNum: 5,
		SecondNum: 10,
	}

	response, err := client.Add(context.Background() ,request)
	if err != nil {
		log.Fatalf("Error getting response for Add function %v", err)
	}

	log.Printf("Output Value is: %v", response.Sum)
}

func doServerStreaming(client calculator.CalculatorServiceClient) {
	request := &calculator.PrimeRequest{
		Num: 120,
	}

	fmt.Println("Sending the request to Server")
	stream, err := client.Prime(context.Background(), request)
	for {
		if err != nil {
			log.Fatalf("error in reading from server %v", err)
		} else {
			resp, err := stream.Recv()
			if err == io.EOF {
				log.Println("Completed all the values")
				break
			} else if err != nil {
				fmt.Printf("Error in getting info from stream %v", err)
			} else {
				fmt.Printf("Got the prime number: %v", resp.GetPrimeNum())
			}
		}
	}
}

func doClientStreaming(client calculator.CalculatorServiceClient) {
	stream, err := client.Average(context.Background())
	if err != nil {
		log.Fatalf("Error connecting to server %v\n", err)
	}
	nums := []int32 {1, 2, 3}
	for _, n := range(nums) {
		fmt.Printf("Sending request %v\n", n)
		stream.Send(&calculator.AverageRequest{
			Num: n,
		})
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error getting response from server %v\n", err)
	}
	fmt.Printf("Response from server %v", res)
}

func doBidirectionalStreaming(client calculator.CalculatorServiceClient) {
	stream, err := client.FindMax(context.Background())
	if err != nil {
		log.Fatalf("Error in getting stream from server %v", err)
	}
	nums := []int32{1, 10, 2, 20, 21, 5, 3, 23}

	waitc := make(chan struct{})

	go func() {
		for _, num := range(nums) {
			err = stream.Send(&calculator.FindMaxRequest{
				Num: num,
			})
			if err != nil {
				log.Fatalf("Error in sending data to server %v", err)
			}
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("All the values are received from the server")
				break
			}
			if err != nil {
				log.Fatalf("error in getting data from server %v", err)
				break
			}
			fmt.Printf("Max Value: %v\n", resp.GetMaxNum())
		}
		close(waitc)
	}()

	<- waitc
}

func doErrorUnary(client calculator.CalculatorServiceClient) {
	fmt.Println("Starting sqrt with error handling")

	doErrHandling(client, 10)
	doErrHandling(client, -2)
}

func doErrHandling(client calculator.CalculatorServiceClient, n int32) {
	resp, err := client.SquareRoot(
		context.Background(),
		&calculator.SquareRootRequest{
			Number: n,
		},
	)
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			fmt.Println("Error Message from Server: " + respErr.Message())
			fmt.Printf("gRPC Error Code: %v\n", respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("Sent a negative number")
				return
			}
		} else {
			log.Fatalf("Some gRPC level error")
		}
	}
	fmt.Printf("Result of square root of %v\n", resp.GetNumberRoot())

}