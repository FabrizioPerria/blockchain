package node

import (
	"testing"

	"github.com/fabrizioperria/blockchain/types"
	"github.com/fabrizioperria/blockchain/utils"
	"github.com/stretchr/testify/assert"
)

func TestAddBlock(t *testing.T) {
	bs := NewMemoryBlockStorer()
	c := NewChain(bs)
	block := utils.GenerateBlock(t, 1)
	assert.NoError(t, c.AddBlock(block))
	hash := types.HashBlockSHA256(block)

	fetchedBlock, err := c.GetBlockByHash(hash)
	assert.NoError(t, err)
	assert.Equal(t, block, fetchedBlock)
}

func TestChainHeight(t *testing.T) {
	bs := NewMemoryBlockStorer()
	c := NewChain(bs)
	genesis, err := c.GetBlockByHeight(0)
	assert.NoError(t, err)
	prev := types.HashBlockSHA256(genesis)
	for i := 1; i < 100; i++ {
		block := utils.GenerateBlock(t, int32(i))
		block.Header.PreviousHash = prev
		assert.NoError(t, c.AddBlock(block))
		assert.Equal(t, int32(i), c.headers.Height())
		prev = types.HashBlockSHA256(block)
	}
}
