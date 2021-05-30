package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"time"

	"github.com/worldofprasanna/grpc-go-code/calculator/calculator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct{
	calculator.UnimplementedCalculatorServiceServer
}

func (*server) Add(ctx context.Context, request *calculator.MathRequest) (*calculator.MathResponse, error) {
	firstNumber := request.FirstNum
	secondNumber := request.SecondNum

	return &calculator.MathResponse{
		Sum: firstNumber + secondNumber,
	}, nil
}

func (*server) Prime(req *calculator.PrimeRequest, stream calculator.CalculatorService_PrimeServer) error {
	num := req.Num
	var i int32 = 2

	fmt.Printf("Got request from client %v\n", num)
	for num > 1 {
		fmt.Printf("Num value: %v, %v\n", num, i)
		if num % i == 0 {
			resp := &calculator.PrimeResponse{
				PrimeNum: i,
			}
			num = num / i
			stream.Send(resp)
			time.Sleep(1000 * time.Millisecond)
		} else {
			i = i + 1
		}
	}

	return nil
}

func (*server) Average(stream calculator.CalculatorService_AverageServer) error {
	fmt.Println("Starting the Average RPC Method")
	sum := int32(0)
	n := int32(0)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			avg := float64(sum / n)
			fmt.Printf("Average computed, %v", avg)
			return stream.SendAndClose(&calculator.AverageResponse{
				Result: avg,
			})
		} else if err != nil {
			log.Fatalf("Error while getting values from Client, %v\n", err)
		} else {
			sum += req.Num
			n++
		}
	}
}

func (*server) FindMax(stream calculator.CalculatorService_FindMaxServer) error {
	currentMax := int32(-1)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("All the messages received from the client")
			break
		} else if err != nil {
			log.Fatalf("Error while getting info from client %v", err)
			return err
		}
		num := req.GetNum()
		if num > currentMax {
			currentMax = num
			fmt.Printf("Max Found: %v, ", num)
			stream.Send(&calculator.FindMaxResponse{
				MaxNum: num,
			})
		}
	}
	return nil
}

func (*server) SquareRoot(ctx context.Context, req *calculator.SquareRootRequest) (*calculator.SquareRootResponse, error) {
	fmt.Println("In Square Root RPC")

	number := req.GetNumber()
	if (number < 0) {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received Invalid Number %v\n", number),
		)
	}
	return &calculator.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

func main() {
	fmt.Println("Hello World from server")

	conn, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		log.Fatalf("Error starting the server %v", err)
	}

	s := grpc.NewServer()
	calculator.RegisterCalculatorServiceServer(s, &server{})

	reflection.Register(s)

	if err = s.Serve(conn); err != nil {
		log.Fatalf("error in starting the grpc server %v", err)
	}
}
