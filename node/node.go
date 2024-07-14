package node

import (
	"context"
	"net"
	"sync"

	"github.com/beevik/guid"
	"github.com/fabrizioperria/blockchain/logging"
	proto "github.com/fabrizioperria/blockchain/protobuf"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type nodeData struct {
	version string
	height  int32
}

type Node struct {
	proto.UnimplementedNodeServer
	peers      map[proto.NodeClient]*proto.HandshakeMsg
	logger     *logrus.Logger
	version    string
	listenAddr string
	id         string
	peerLock   sync.RWMutex
	height     int32
}

func New() *Node {
	d := getNodeData()
	id := guid.New().String()

	return &Node{
		version:  d.version,
		height:   d.height,
		peers:    make(map[proto.NodeClient]*proto.HandshakeMsg),
		peerLock: sync.RWMutex{},
		id:       id,
		logger:   logging.LoggerFactory("logs/" + id + ".log"),
	}
}

func getNodeData() *nodeData {
	return &nodeData{
		version: "1.0.0",
		height:  100,
	}
}

func (n *Node) Start(listenAddr string) {
	n.listenAddr = listenAddr
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		n.logger.Fatalf("failed to listen: %v", err)
	}

	proto.RegisterNodeServer(grpcServer, n)

	n.logger.Info("Server started on port 3000")

	grpcServer.Serve(listener)
}

func (n *Node) BootstrapConnect(addresses []string) error {
	for _, address := range addresses {
		client, err := makeNodeClient(address)
		if err != nil {
			n.logger.Fatalf("failed to dial server: %v", err)
			return err
		}

		msg, err := client.Handshake(context.Background(), &proto.HandshakeMsg{
			Version: n.version,
			Height:  n.height,
			Address: n.listenAddr,
		})
		if err != nil {
			n.logger.Fatalf("failed to make handshake: %v", err)
			continue
		}

		n.addPeer(&client, msg)
	}

	return nil
}

func (n *Node) HandleTransaction(ctx context.Context, transaction *proto.Transaction) (*proto.Ack, error) {
	n.logger.Info("Received transaction")
	n.logger.WithFields(logrus.Fields{
		"version": transaction.Version,
	}).Info("Transaction received")
	return &proto.Ack{}, nil
}

func (n *Node) Handshake(ctx context.Context, helo *proto.HandshakeMsg) (*proto.HandshakeMsg, error) {
	n.logger.WithFields(logrus.Fields{
		"fromAddress": n.listenAddr,
		"toAddress":   helo.Address,
		"version":     helo.Version,
		"height":      helo.Height,
	}).Info("Received handshake")

	client, err := makeNodeClient(helo.Address)
	if err != nil {
		n.logger.Fatalf("failed to dial server: %v", err)
		return nil, err
	}

	myMsg := &proto.HandshakeMsg{
		Version: n.version,
		Height:  n.height,
		Address: n.listenAddr,
	}
	n.addPeer(&client, helo)

	return myMsg, nil
}

func (n *Node) addPeer(peer *proto.NodeClient, data *proto.HandshakeMsg) bool {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()
	n.peers[*peer] = data
	n.logger.WithFields(logrus.Fields{
		"fromAddress": n.listenAddr,
		"toAddress":   data.Address,
		"version":     data.Version,
		"height":      data.Height,
	}).Info("Added peer")
	return true
}

func (n *Node) removePeer(peer *proto.NodeClient) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()
	delete(n.peers, *peer)
	n.logger.WithFields(logrus.Fields{
		"address": *peer,
	}).Info("Removed peer")
}

func makeNodeClient(listenAddr string) (proto.NodeClient, error) {
	client, err := grpc.NewClient(listenAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return proto.NewNodeClient(client), nil
}
