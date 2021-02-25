package main

import (
	"context"
	"dumanModule/greet/greetpb"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	fmt.Println("Hello I'm a client")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := greetpb.NewCalculatorServiceClient(conn)

	doPrimeNumberDecomposition(c)

	//doComputeAverage(c)
}

func doPrimeNumberDecomposition(c greetpb.CalculatorServiceClient) {
	ctx := context.Background()
	req := &greetpb.CalculatorRequest{Calculator: &greetpb.Calculator{
		Num: 120,
	}}
	log.Printf("sending: %v \n", req.Calculator.GetNum())
	stream, err := c.PrimeNumberDecomposition(ctx, req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC %v", err)
	}
	defer stream.CloseSend()
LOOP:
	for {
		res, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break LOOP
			}
			log.Fatalf("error while reciving from PrimeNumberDecomposition RPC %v", err)
		}
		log.Printf("response from DecomposePrimeNumber:%v \n", res.GetResult())
	}
}

func doComputeAverage(c greetpb.CalculatorServiceClient) {
	requests := []*greetpb.CalculatorRequest{
		{
			Calculator: &greetpb.Calculator{
				Num: 1,
			},
		},
		{
			Calculator: &greetpb.Calculator{
				Num: 2,
			},
		},
		{
			Calculator: &greetpb.Calculator{
				Num: 3,
			},
		},
		{
			Calculator: &greetpb.Calculator{
				Num: 4,
			},
		},
	}
	ctx := context.Background()
	stream, err := c.ComputeAverage(ctx)
	if err != nil {
		log.Fatalf("error while calling ComputeAverage: %v", err)
	}
	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(time.Second)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from ComputeAverage: %v", err)
	}
	fmt.Printf("ComputeAverage Response: %v\n", res)
}
