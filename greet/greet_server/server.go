package main

import (
	"dumanModule/greet/greetpb"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

//Server with embedded UnimplementedGreetServiceServer
type Server struct {
	greetpb.UnimplementedCalculatorServiceServer
}

// PrimeNumberDecomposition is an example of stream from server side
func (s *Server) PrimeNumberDecomposition(req *greetpb.CalculatorRequest, stream greetpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("DecomposePrimeNumber function was invoked with %v \n", req)
	num := req.GetCalculator().GetNum()
	var i int32

	for i = 2; i <= num; i++ {
		for num%i == 0 {
			res := &greetpb.CalculatorResponse{Result: float32(i)}
			if err := stream.Send(res); err != nil {
				log.Fatalf("error while sending DecomposePrimeNumber responses: %v", err.Error())
			}
			num /= i
			time.Sleep(time.Second)
		}
	}
	return nil
}

// ComputeAverage
func (s *Server) ComputeAverage(stream greetpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("LongGreet function was invoked with a streaming request\n")
	var sum int32
	var counter int32 = 0
	var result float32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// we have finished reading the client stream
			result = float32(sum) / float32(counter)
			return stream.SendAndClose(&greetpb.CalculatorResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		num := req.Calculator.GetNum()
		sum += num
		counter++
	}
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen:%v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterCalculatorServiceServer(s, &Server{})
	log.Println("Server is running on port:50051")
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve:%v", err)
	}
}
