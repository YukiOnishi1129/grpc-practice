package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"pancake.maker/get/api"
	"pancake.maker/handler"
)

func main() {
	port := 50051
	lis, err := net.Listen("tcp",fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: &v", err)
	}
	server := grpc.NewServer()
	// .protoに定義したPancakeBakerServiceに対応している
	// 第二引数のハンドラに対応するメソッド(bakeかreport)を呼び出す
	api.RegisterPancakeBakerServiceServer(server, handler.NewBakerHandler())
	reflection.Register(server)

	go func() {
		log.Printf("start gRPC server port: &v", port)
		server.Serve(lis)
	} ()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	server.GracefulStop()
}