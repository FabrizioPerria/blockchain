package types

import (
	"crypto/sha256"

	crypto "github.com/fabrizioperria/blockchain/crypto"
	proto "github.com/fabrizioperria/blockchain/protobuf"
	pb "google.golang.org/protobuf/proto"
)

func SignTransaction(transaction *proto.Transaction, privateKey *crypto.PrivateKey) *crypto.Signature {
	hash := HashTransactionSHA256(transaction)
	signature := privateKey.Sign(hash)

	return signature
}

func HashTransactionSHA256(transaction *proto.Transaction) []byte {
	b, err := pb.Marshal(transaction)
	if err != nil {
		panic(err)
	}

	hash := sha256.Sum256(b)
	return hash[:]
}

func VerifyTransaction(transaction *proto.Transaction) bool {
	for _, input := range transaction.Inputs {
		signature := crypto.SignatureFromBytes(input.Signature)
		publicKey := crypto.PublicKeyFromBytes(input.PublicKey)

		hash := HashTransactionSHA256(transaction)
		if !signature.Verify(publicKey, hash) {
			return false
		}
	}
	return true
}
