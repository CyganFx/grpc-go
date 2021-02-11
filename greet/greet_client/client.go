package main

import (
	"context"
	"dumanModule/greet/greetpb"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm a client")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)

	doGreetingMax(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	ctx := context.Background()
	request := &greetpb.GreetRequest{Greeting: &greetpb.Greeting{
		FirstName: "Duman",
		LastName:  "Ishanov",
	}}

	response, err := c.Greet(ctx, request)
	if err != nil {
		log.Fatalf("error while calling Greet RPC %v", err)
	}
	log.Printf("response from Greet:%v", response.Result)
}

func doManyTimesFromServer(c greetpb.GreetServiceClient) {
	ctx := context.Background()
	req := &greetpb.GreetManyTimesRequest{Greeting: &greetpb.Greeting{
		FirstName: "Tleu",
		LastName:  "Mukatayev",
	}}

	stream, err := c.GreetManyTimes(ctx, req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC %v", err)
	}
	defer stream.CloseSend()

LOOP:
	for {
		res, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				// we've reached the end of the stream
				break LOOP
			}
			log.Fatalf("error while reciving from GreetManyTimes RPC %v", err)
		}
		log.Printf("response from GreetManyTimes:%v \n", res.GetResult())
	}

}

func doLongGreet(c greetpb.GreetServiceClient) {

	requests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Tleu",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Bob",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Alice",
			},
		},
	}

	ctx := context.Background()
	stream, err := c.LongGreet(ctx)
	if err != nil {
		log.Fatalf("error while calling LongGreet: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from LongGreet: %v", err)
	}
	fmt.Printf("LongGreet Response: %v\n", res)
}

func doGreetingEveryone(c greetpb.GreetServiceClient) {
	log.Printf("Starting to do a BIDI Streaming RPC...")
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
	}

	waitc := make(chan struct{})

	requests := []*greetpb.GreetEveryoneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Tleu",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Bob",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Alice",
			},
		},
	}

	go func() {
		for _, req := range requests {
			log.Printf("sending message: %v", req)
			err := stream.Send(req)
			if err != nil {
				log.Fatalf("error while sending req to server: %v", err)
			}
			time.Sleep(time.Second)
		}
		err := stream.CloseSend()
		if err != nil {
			log.Fatalf("error while closing client's stream")
		}

	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while receiving getting message from stream %v", err)
			}
			log.Printf("Received: %v", res.GetResult())
		}
		close(waitc)
	}()
	<-waitc
}

func doGreetingMax(c greetpb.GreetServiceClient) {

	log.Printf("Starting to do a BIDI Streaming RPC...")
	stream, err := c.GreetMax(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
	}
	waitc := make(chan struct{})

	requests := []*greetpb.GreetMaxRequest{
		{
			GreetingMax: &greetpb.GreetingMax{
				Num: 1,
			},
		},
		{
			GreetingMax: &greetpb.GreetingMax{
				Num: 5,
			},
		},
		{
			GreetingMax: &greetpb.GreetingMax{
				Num: 3,
			},
		},
		{
			GreetingMax: &greetpb.GreetingMax{
				Num: 20,
			},
		},
	}

	go func() {
		for _, req := range requests {
			log.Printf("sending message: %v", req)

			err := stream.Send(req)
			if err != nil {
				log.Fatalf("error while sending req to server: %v", err)
			}
			time.Sleep(time.Second)
		}
		err := stream.CloseSend()
		if err != nil {
			log.Fatalf("error while closing client's stream")
		}
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while receiving getting message from stream %v", err)
			}
			log.Printf("Current local max is:: %v", res.GetResult())
		}
		close(waitc)
	}()
	<-waitc
}
