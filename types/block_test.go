package types

import (
	"testing"

	"github.com/fabrizioperria/blockchain/crypto"
	"github.com/fabrizioperria/blockchain/utils"
	"github.com/stretchr/testify/assert"
)

func TestHashBlockSHA256(t *testing.T) {
	block := utils.GenerateBlock(t, 1)

	hash := HashBlockSHA256(block)
	assert.Equal(t, 32, len(hash))
}

func TestSignBlock(t *testing.T) {
	block := utils.GenerateBlock(t, 1)
	privateKey := crypto.GeneratePrivateKey()
	publicKey := privateKey.Public()

	signature := SignBlock(block, privateKey)
	assert.NotNil(t, signature)
	assert.Equal(t, 64, len(signature.Bytes()))
	assert.True(t, signature.Verify(publicKey, HashBlockSHA256(block)))
}
