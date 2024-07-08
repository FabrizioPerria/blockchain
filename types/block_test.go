package types

import (
	crand "crypto/rand"
	mrand "math/rand"
	"testing"
	"time"

	"github.com/fabrizioperria/blockchain/crypto"
	proto "github.com/fabrizioperria/blockchain/protobuf"
	"github.com/stretchr/testify/assert"
)

func RandomHash(t *testing.T) []byte {
	hash := make([]byte, 32)
	n, err := crand.Read(hash)
	assert.NoError(t, err)
	assert.NotZero(t, n)

	return hash
}

func GenerateBlock(t *testing.T) *proto.Block {
	return &proto.Block{
		Header: &proto.Header{
			Version:      1,
			Height:       int32(mrand.Intn(1000) + 1),
			PreviousHash: RandomHash(t),
			MerkleRoot:   RandomHash(t),
			Timestamp:    time.Now().UnixNano(),
		},
	}
}

func TestHashBlockSHA256(t *testing.T) {
	block := GenerateBlock(t)

	hash := HashBlockSHA256(block)
	assert.Equal(t, 32, len(hash))
}

func TestSignBlock(t *testing.T) {
	block := GenerateBlock(t)
	privateKey := crypto.GeneratePrivateKey()
	publicKey := privateKey.Public()

	signature := SignBlock(block, privateKey)
	assert.NotNil(t, signature)
	assert.Equal(t, 64, len(signature.Bytes()))
	assert.True(t, signature.Verify(publicKey, HashBlockSHA256(block)))
}
