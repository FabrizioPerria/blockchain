package types

import (
	"crypto/sha256"

	crypto "github.com/fabrizioperria/blockchain/crypto"
	proto "github.com/fabrizioperria/blockchain/protobuf"
	pb "google.golang.org/protobuf/proto"
)

func HashBlockSHA256(block *proto.Block) []byte {
	return HashHeaderSHA256(block.Header)
}

func SignBlock(block *proto.Block, privateKey *crypto.PrivateKey) *crypto.Signature {
	hash := HashBlockSHA256(block)
	return privateKey.Sign(hash)
}

func HashHeaderSHA256(header *proto.Header) []byte {
	b, err := pb.Marshal(header)
	if err != nil {
		panic(err)
	}

	hash := sha256.Sum256(b)
	return hash[:]
}
