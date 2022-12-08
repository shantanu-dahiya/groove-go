package crypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"

	parser "groove-go.nyu.edu/pkg/parser"
)

var ec = elliptic.P256()

type PublicKey ecdsa.PublicKey
type PrivateKey []byte

func GenerateKeys() (PublicKey, PrivateKey, error) {
	publicKey := PublicKey{}
	privateKey, x, y, err := elliptic.GenerateKey(ec, rand.Reader)
	if err != nil {
		return publicKey, privateKey, err
	}

	publicKey.X = x
	publicKey.Y = y
	publicKey.Curve = ec

	return publicKey, privateKey, nil
}

// taken from ECDH
func ComputeSymmetricKey(peerPublicKey PublicKey, privateKey PrivateKey) []byte {
	sX, _ := ec.ScalarMult(peerPublicKey.X, peerPublicKey.Y, privateKey)
	h := sha256.New()
	h.Write(sX.Bytes())
	hashedKey := h.Sum(nil)

	return hashedKey
}

func MarshalPublicKey(publicKey PublicKey) []byte {
	return elliptic.Marshal(ec, publicKey.X, publicKey.Y)
}

func UnmarshalPublicKey(data []byte) PublicKey {
	x, y := elliptic.Unmarshal(ec, data)
	return PublicKey{ec, x, y}
}

func EncryptSymmetric(message []byte, key []byte) ([]byte, error) {
	message = PKCS5Padding(message, aes.BlockSize)
	aes_cipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	IV := bytes.Repeat([]byte{byte(0)}, aes_cipher.BlockSize()) // TODO: make this random

	ciphertext := make([]byte, len(message))
	aes_cipher.Encrypt(ciphertext, []byte(message))
	mode := cipher.NewCBCEncrypter(aes_cipher, IV)
	mode.CryptBlocks(ciphertext, message)

	return ciphertext, nil
}

func DecryptSymmetric(data []byte, key []byte) ([]byte, error) {
	aes_cipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	IV := bytes.Repeat([]byte{byte(0)}, aes_cipher.BlockSize()) // TODO: make this random

	mode := cipher.NewCBCDecrypter(aes_cipher, IV)
	mode.CryptBlocks(data, data)

	return parser.RemoveTrailingSpaces(data), err
}
