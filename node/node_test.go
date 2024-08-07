package node

import (
	"strconv"
	"testing"
	"time"

	// "time"

	"github.com/stretchr/testify/assert"
)

func TestSetupCluster(t *testing.T) {
	n := []*Node{}
	n = append(n, makeNode("localhost:3000", []string{}))
	expectedNumPeers := 9
	for i := 0; i < expectedNumPeers; i++ {
		port := 3001 + i
		n = append(n, makeNode("localhost:"+strconv.Itoa(port), []string{"localhost:3000"}))
	}

	time.Sleep(1 * time.Second)
	for _, node := range n {
		l := len(node.GetPeers())
		assert.Equal(t, expectedNumPeers, l, "expected %d peers, got %d", expectedNumPeers, l)
	}
}

func makeNode(listenAddr string, bootstrapNodes []string) *Node {
	n := New()
	go n.Start(listenAddr, bootstrapNodes)
	time.Sleep(1 * time.Second)

	return n
}
