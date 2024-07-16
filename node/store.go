package node

import (
	"encoding/hex"
	"fmt"
	"sync"

	proto "github.com/fabrizioperria/blockchain/protobuf"
	"github.com/fabrizioperria/blockchain/types"
)

type BlockStorer interface {
	Put(*proto.Block) error
	Get(string) (*proto.Block, error)
}

type MemoryBlockStorer struct {
	blocks sync.Map
}

func NewMemoryBlockStorer() *MemoryBlockStorer {
	return &MemoryBlockStorer{blocks: sync.Map{}}
}

func (m *MemoryBlockStorer) Put(block *proto.Block) error {
	hash := hex.EncodeToString(types.HashBlockSHA256(block))
	m.blocks.Store(hash, block)
	return nil
}

func (m *MemoryBlockStorer) Get(hash string) (*proto.Block, error) {
	block, ok := m.blocks.Load(hash)
	if !ok {
		return nil, fmt.Errorf("block %s not found", hash)
	}

	return block.(*proto.Block), nil
}
