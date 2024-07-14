package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

const (
	signatureSize  = ed25519.SignatureSize
	privateKeySize = ed25519.PrivateKeySize
	publicKeySize  = ed25519.PublicKeySize
	seedSize       = ed25519.SeedSize
	addressSize    = 20
)

type PrivateKey struct {
	key ed25519.PrivateKey
}

func GeneratePrivateKey() *PrivateKey {
	b := make([]byte, seedSize)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		// If you can't generate a seed, just panic, the application wouldn't work anyways as the key generation will panic.
		panic(err)
	}

	return GeneratePrivateKeyFromSeed(b)
}

func GeneratePrivateKeyFromSeed(seed []byte) *PrivateKey {
	if len(seed) != seedSize {
		panic("invalid seed size")
	}
	pk := ed25519.NewKeyFromSeed(seed)
	return &PrivateKey{key: pk}
}

func (p *PrivateKey) Bytes() []byte {
	return p.key
}

func (p *PrivateKey) Sign(data []byte) *Signature {
	return &Signature{
		data: ed25519.Sign(p.Bytes(), data),
	}
}

func (p *PrivateKey) Public() *PublicKey {
	return &PublicKey{
		key: p.key.Public().(ed25519.PublicKey),
	}
}

// ====================================================================================================

type PublicKey struct {
	key ed25519.PublicKey
}

func (p *PublicKey) Bytes() []byte {
	return p.key
}

func PublicKeyFromBytes(data []byte) *PublicKey {
	if len(data) != publicKeySize {
		panic(fmt.Errorf(`invalid public key size. Size must be %d, but got %d`, publicKeySize, len(data)))
	}

	return &PublicKey{
		key: ed25519.PublicKey(data),
	}
}

// ====================================================================================================

type Signature struct {
	data []byte
}

func (s *Signature) Bytes() []byte {
	return s.data
}

func SignatureFromBytes(data []byte) *Signature {
	if len(data) != signatureSize {
		panic(fmt.Errorf(`invalid signature size. Size must be %d, but got %d`, signatureSize, len(data)))
	}
	return &Signature{
		data: data,
	}
}

func (s *Signature) Verify(pubKey *PublicKey, data []byte) bool {
	return ed25519.Verify(pubKey.Bytes(), data, s.Bytes())
}

// ====================================================================================================

type Address struct {
	data []byte
}

func (p *PublicKey) Address() *Address {
	return &Address{
		data: p.Bytes()[:addressSize],
	}
}

func (a *Address) Bytes() []byte {
	return a.data
}

func (a *Address) String() string {
	return hex.EncodeToString(a.Bytes())
}
