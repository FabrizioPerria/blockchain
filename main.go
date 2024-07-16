package main

import (
	"context"
	"strconv"
	"time"

	"github.com/fabrizioperria/blockchain/logging"
	"github.com/fabrizioperria/blockchain/node"
	proto "github.com/fabrizioperria/blockchain/protobuf"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	n := []*node.Node{}
	n = append(n, makeNode("localhost:3000", []string{}))
	for i := 1; i < 30; i++ {
		port := 3000 + i
		n = append(n, makeNode("localhost:"+strconv.Itoa(port), []string{"localhost:3000"}))
	}

	// expectedNumPeers := 9
	// for i, node := range n {
	// 	l := len(node.GetPeers())
	// 	if l != expectedNumPeers {
	// 		log.Fatalf("[%d] expected %d peers, got %d", i, expectedNumPeers, l)
	// 	}
	// }

	time.Sleep(2 * time.Second)
	for i, node := range n {
		log.Infof("[%d] peers: %v", i, node.GetPeers())
	}
	// select {}
}

var log = logging.LoggerFactory("logs/log.log")

func makeNode(listenAddr string, bootstrapNodes []string) *node.Node {
	n := node.New()
	go n.Start(listenAddr, bootstrapNodes)
	time.Sleep(1 * time.Second)

	return n
}

func makeTransaction(clientAddr string) {
	client, err := grpc.NewClient(clientAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
