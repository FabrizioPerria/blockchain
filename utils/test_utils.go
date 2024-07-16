package utils

import (
	"testing"
	"time"

	crand "crypto/rand"

	"github.com/stretchr/testify/assert"

	proto "github.com/fabrizioperria/blockchain/protobuf"
)

func RandomHash(t *testing.T) []byte {
	hash := make([]byte, 32)
	n, err := crand.Read(hash)
	assert.NoError(t, err)
	assert.NotZero(t, n)

	return hash
}

func GenerateBlock(t *testing.T, height int32) *proto.Block {
	return &proto.Block{
		Header: &proto.Header{
			Version:      1,
			Height:       height,
			PreviousHash: RandomHash(t),
			MerkleRoot:   RandomHash(t),
			Timestamp:    time.Now().UnixNano(),
		},
	}
}
