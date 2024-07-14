package types

import (
	"testing"

	"github.com/fabrizioperria/blockchain/crypto"
	proto "github.com/fabrizioperria/blockchain/protobuf"
	"github.com/stretchr/testify/assert"
)

func TestNewTransaction(t *testing.T) {
	fromPrivateKey := crypto.GeneratePrivateKey()
	fromPublicKey := fromPrivateKey.Public()
	fromAddress := fromPublicKey.Address()

	toPrivateKey := crypto.GeneratePrivateKey()
	toPublicKey := toPrivateKey.Public()
	toAddress := toPublicKey.Address()

	totalAmount := int64(800)
	amount := int64(100)

	input := &proto.TxInput{
		PreviousTxHash:  []byte(""),
		PrevOutputIndex: 0,
		PublicKey:       fromPublicKey.Bytes(),
	}

	output1 := &proto.TxOutput{
		Amount:      amount,
		DestAddress: toAddress.Bytes(),
	}

	output2 := &proto.TxOutput{
		Amount:      totalAmount - amount,
		DestAddress: fromAddress.Bytes(),
	}

	transaction := &proto.Transaction{
		Version: 1,
		Inputs:  []*proto.TxInput{input},
		Outputs: []*proto.TxOutput{output1, output2},
	}

	signature := SignTransaction(transaction, fromPrivateKey)
	input.Signature = signature.Bytes()

	assert.True(t, VerifyTransaction(transaction))
}
