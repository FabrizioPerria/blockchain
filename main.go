package main

import (
	"context"
	// "time"

	"github.com/fabrizioperria/blockchain/logging"
	"github.com/fabrizioperria/blockchain/node"
	proto "github.com/fabrizioperria/blockchain/protobuf"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	makeNode("localhost:3000", []string{})
	makeNode("localhost:4000", []string{"localhost:3000"})

	// go func() {
	// 	for {
	// 		time.Sleep(2 * time.Second)
	// 		makeTransaction()
	// 	}
	// }()

	select {}
}

var log = logging.LoggerFactory("logs/log.log")

func makeNode(listenAddr string, bootstrapNodes []string) *node.Node {
	n := node.New()
	go n.Start(listenAddr)
	if len(bootstrapNodes) > 0 {
		if err := n.BootstrapConnect(bootstrapNodes); err != nil {
			log.Fatalf("failed to connect to bootstrap nodes: %v", err)
		}
	}

	return n
}

func makeTransaction() {
	client, err := grpc.NewClient("localhost:3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial server: %v", err)
	}
	defer client.Close()

	nodeClient := proto.NewNodeClient(client)

	transaction := &proto.Transaction{
		Version: 1,
		Inputs:  []*proto.TxInput{},
		Outputs: []*proto.TxOutput{},
	}

	_, err = nodeClient.Handshake(context.Background(), &proto.HandshakeMsg{
		Version: "1.0.0",
		Height:  100,
		Address: ":4000",
	})
	if err != nil {
		log.Fatalf("failed to make handshake: %v", err)
	}

	_, err = nodeClient.HandleTransaction(context.Background(), transaction)
	if err != nil {
		log.Fatalf("failed to make transaction: %v", err)
	}

	log.Info("Transaction made successfully")
}
