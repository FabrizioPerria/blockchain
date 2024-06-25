package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getStaticSeed() []byte {
	return []byte{
		0x11, 0xA2, 0x93, 0xF4, 0x11, 0xA2, 0x93, 0xF4,
		0x11, 0xA2, 0x93, 0xF4, 0x11, 0xA2, 0x93, 0xF4,
		0x11, 0xA2, 0x93, 0xF4, 0x11, 0xA2, 0x93, 0xF4,
		0x11, 0xA2, 0x93, 0xF4, 0x11, 0xA2, 0x93, 0xF4,
	}
}

func getStaticPrivateKey() *PrivateKey {
	return GeneratePrivateKeyFromSeed(getStaticSeed())
}

func TestGenerateNewPrivateKey(t *testing.T) {
	pk := GeneratePrivateKey()

	assert.NotNil(t, pk)
	assert.Equal(t, len(pk.Bytes()), privateKeySize)
}

func TestGeneratePrivateKeyFromSeedDeterministically(t *testing.T) {
	pk := GeneratePrivateKeyFromSeed(getStaticSeed())

	assert.NotNil(t, pk)
	assert.Equal(t, len(pk.Bytes()), privateKeySize)

	pk2 := GeneratePrivateKeyFromSeed(getStaticSeed())

	assert.Equal(t, pk.Bytes(), pk2.Bytes())
}

func TestPublicKeyFromPrivateKey(t *testing.T) {
	pk := getStaticPrivateKey()
	pub := pk.Public()

	expectedPubKey := []byte{
		0x5e, 0x45, 0x18, 0x16, 0x25, 0xb3, 0x27, 0xdc,
		0x49, 0x8b, 0xe8, 0xef, 0x8c, 0x40, 0xd4, 0xac,
		0xda, 0xc9, 0x81, 0xdc, 0x5, 0x87, 0x4d, 0xbe,
		0x86, 0x92, 0x97, 0xad, 0x4, 0xec, 0x11, 0x72,
	}

	assert.NotNil(t, pub)
	assert.Equal(t, len(pub.Bytes()), publicKeySize)
	assert.Equal(t, pub.Bytes(), expectedPubKey)
}

func TestSignAndVerify(t *testing.T) {
	pk := GeneratePrivateKey()
	pub := pk.Public()

	data := []byte("hello world")
	signature := pk.Sign(data)

	assert.True(t, signature.Verify(pub, data))
}

func TestSignAndVerifyFail(t *testing.T) {
	pk := GeneratePrivateKey()
	pub := pk.Public()

	data := []byte("hello world")
	wrongData := []byte("hello world!")

	signature := pk.Sign(data)

	assert.False(t, signature.Verify(pub, wrongData))
}

func TestSignAndVerifyFailWithDifferentPublicKey(t *testing.T) {
	pk1 := GeneratePrivateKey()
	pk2 := GeneratePrivateKey()

	data := []byte("hello world")
	signature := pk1.Sign(data)

	assert.False(t, signature.Verify(pk2.Public(), data))
}

func TestAddressFromPublicKey(t *testing.T) {
	pk := getStaticPrivateKey()
	pub := pk.Public()

	addr := pub.Address()

	expectedAddress := []byte{
		0x5e, 0x45, 0x18, 0x16, 0x25, 0xb3, 0x27, 0xdc,
		0x49, 0x8b, 0xe8, 0xef, 0x8c, 0x40, 0xd4, 0xac,
		0xda, 0xc9, 0x81, 0xdc,
	}

	assert.NotNil(t, addr)
	assert.Equal(t, len(addr.Bytes()), addressSize)
	assert.Equal(t, addr.Bytes(), expectedAddress)
}

func TestAddressFromPublicKeyIsDeterministic(t *testing.T) {
	pk := GeneratePrivateKey()
	pub := pk.Public()

	addr1 := pub.Address()
	addr2 := pub.Address()

	assert.Equal(t, addr1.Bytes(), addr2.Bytes())
}

func TestAddressFromPublicKeyIsDifferentForDifferentPublicKeys(t *testing.T) {
	pk1 := GeneratePrivateKey()
	pub1 := pk1.Public()

	pk2 := GeneratePrivateKey()
	pub2 := pk2.Public()

	addr1 := pub1.Address()
	addr2 := pub2.Address()

	assert.NotEqual(t, addr1.Bytes(), addr2.Bytes())
}

func TestAddressString(t *testing.T) {
	addr := Address{
		data: []byte{0x11, 0xA2, 0x93, 0xF4},
	}

	assert.Equal(t, "11a293f4", addr.String())
}
