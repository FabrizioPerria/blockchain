package node

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"

	// "sync"

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

type addPeerData struct {
	client *proto.NodeClient
	data   *proto.HandshakeMsg
}

type Node struct {
	proto.UnimplementedNodeServer
	peers        sync.Map
	logger       *logrus.Logger
	addPeerCh    chan *addPeerData
	removePeerCh chan string
	getPeersCh   chan chan []string
	version      string
	listenAddr   string
	id           string
	height       int32
}

func (n *Node) managePeers() {
	for {
		select {
		case peer := <-n.removePeerCh:
			n.peers.Delete(peer)
		case data := <-n.addPeerCh:
			if data.data.Address != "" {
				n.peers.Store(data.data.Address, data)
			}
		case res := <-n.getPeersCh:
			peers := []string{}
			n.peers.Range(func(key, value interface{}) bool {
				peers = append(peers, key.(string))
				return true
			})
			res <- peers
		}
	}
}

func New() *Node {
	d := getNodeData()

	n := &Node{
		version:      d.version,
		height:       d.height,
		peers:        sync.Map{},
		addPeerCh:    make(chan *addPeerData, 100),
		removePeerCh: make(chan string, 100),
		getPeersCh:   make(chan chan []string, 100),
	}
	go n.managePeers()

	return n
}

func getNodeData() *nodeData {
	return &nodeData{
		version: "1.0.0",
		height:  100,
	}
}

func (n *Node) Start(listenAddr string, bootstrapNodes []string) {
	n.listenAddr = listenAddr
	n.logger = logging.LoggerFactory("logs/" + strings.ReplaceAll(listenAddr, ":", "") + ".log")
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		n.logger.Fatalf("failed to listen: %v", err)
	}

	proto.RegisterNodeServer(grpcServer, n)

	if err := n.bootstrapConnect(bootstrapNodes); err != nil {
		n.logger.Fatalf("failed to connect to bootstrap nodes: %v", err)
	}

	n.logger.Infof("Server started on %s", listenAddr)

	grpcServer.Serve(listener)
}

func (n *Node) bootstrapConnect(addresses []string) error {
	n.logger.Infof("[%s] Bootstrapping to %v", n.listenAddr, addresses)
	for _, address := range addresses {
		if n.hasConnectedTo(address) {
			continue
		}

		n.logger.WithFields(logrus.Fields{
			"from":    n.listenAddr,
			"address": address,
		}).Info("Bootstrapping to address")
		client, msg, err := n.dialRemote(address)
		if err != nil {
			return err
		}
		n.addPeer(client, msg)
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
	if n.hasConnectedTo(helo.Address) {
		return nil, nil
	}
	client, err := makeNodeClient(helo.Address)
	if err != nil {
		n.logger.Fatalf("failed to dial server: %v", err)
		return nil, err
	}

	myMsg := &proto.HandshakeMsg{
		Version:    n.version,
		Height:     n.height,
		Address:    n.listenAddr,
		KnownPeers: n.GetPeers(),
	}
	n.addPeer(&client, helo)

	return myMsg, nil
}

func (n *Node) GetPeers() []string {
	res := make(chan []string)
	go func() { n.getPeersCh <- res }()
	return <-res
}

func (n *Node) addPeer(peer *proto.NodeClient, data *proto.HandshakeMsg) bool {
	if n.hasConnectedTo(data.Address) {
		return false
	}
	n.addPeerCh <- &addPeerData{client: peer, data: data}
	n.logger.WithFields(logrus.Fields{
		"receiver":     n.listenAddr,
		"addedPeer":    data.Address,
		"theirVersion": data.Version,
		"theirHeight":  data.Height,
	}).Info("Added peer")

	go n.bootstrapConnect(data.KnownPeers)

	return true
}

func (n *Node) dialRemote(address string) (*proto.NodeClient, *proto.HandshakeMsg, error) {
	if address == n.listenAddr {
		return nil, nil, fmt.Errorf("cannot connect to self")
	}
	client, err := makeNodeClient(address)
	if err != nil {
		return nil, nil, err
	}
	msg, err := client.Handshake(context.Background(), &proto.HandshakeMsg{
		Version:    n.version,
		Height:     n.height,
		Address:    n.listenAddr,
		KnownPeers: n.GetPeers(),
	})
	if err != nil {
		return nil, nil, err
	}
	return &client, msg, nil
}

func (n *Node) hasConnectedTo(address string) bool {
	if address == n.listenAddr {
		return true
	}

	_, ok := n.peers.Load(address)
	return ok
}

func (n *Node) removePeer(peer string) {
	n.removePeerCh <- peer
	n.logger.WithFields(logrus.Fields{
		"address": peer,
	}).Info("Removed peer")
}

func makeNodeClient(listenAddr string) (proto.NodeClient, error) {
	client, err := grpc.NewClient(listenAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return proto.NewNodeClient(client), nil
}
