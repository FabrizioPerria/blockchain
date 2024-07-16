package node

import (
	"encoding/hex"
	"fmt"

	proto "github.com/fabrizioperria/blockchain/protobuf"
	"github.com/fabrizioperria/blockchain/types"
)

type Headerer interface {
	Length() int32
	Height() int32
	Add(*proto.Header) error
}

type HeadersChain struct {
	Headerer
	headers []*proto.Header
}

func (hc *HeadersChain) Height() int32 {
	return hc.Length() - 1
}

func (hc *HeadersChain) Length() int32 {
	return int32(len(hc.headers))
}

func (hc *HeadersChain) Add(header *proto.Header) error {
	hc.headers = append(hc.headers, header)
	return nil
}

type Chain struct {
	blockStorer BlockStorer
	headers     *HeadersChain
}

func NewChain(blockStorer BlockStorer) *Chain {
	chain := &Chain{
		blockStorer: blockStorer,
		headers:     &HeadersChain{headers: []*proto.Header{}},
	}

	genesisBlock := &proto.Block{
		Header: &proto.Header{
			Version:      1,
			Height:       0,
			PreviousHash: make([]byte, 32),
			MerkleRoot:   make([]byte, 32),
			Timestamp:    0,
		},
	}
	chain.AddBlock(genesisBlock)

	return chain
}

func (c *Chain) AddBlock(block *proto.Block) error {
	c.headers.Add(block.Header)
	return c.blockStorer.Put(block)
}

func (c *Chain) GetBlockByHash(hash []byte) (*proto.Block, error) {
	h := hex.EncodeToString(hash)
	return c.blockStorer.Get(h)
}

func (c *Chain) GetBlockByHeight(height int32) (*proto.Block, error) {
	if height < 0 || height >= c.headers.Length() {
		return nil, fmt.Errorf("block height %d out of range", height)
	}
	header := c.headers.headers[height]
	hash := types.HashHeaderSHA256(header)
	return c.GetBlockByHash(hash)
}

func (c *Chain) Height() int32 {
	return c.headers.Height()
}
