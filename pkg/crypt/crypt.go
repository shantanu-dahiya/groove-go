package crypt

import (
	"crypto/aes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"

	"groove-go.nyu.edu/pkg/global"
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
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(message))
	cipher.Encrypt(out, []byte(message))
	return out, nil
}

func DecryptSymmetric(data []byte, key []byte) ([]byte, error) {
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(data))
	cipher.Decrypt(out, []byte(data))
	return out, nil
}

func EncryptOnion(message []byte, tagId int, publicKey PublicKey, ckt global.Circuit) ([]byte, error) {
	for i, e := range ckt {
		message, err := EncryptSymmetric(message, e.SymmetricKey)
		var nextHopAddress string
		var nextHopPort int = 0
		if i < len(ckt)-1 {
			nextHopAddress = ckt[i+1].Server.Host
			nextHopPort = ckt[i+1].Server.Port
		}
		cfd := &global.CircuitForwardingData{
			Data:           message,
			NextHopAddress: nextHopAddress,
			NextHopPort:    nextHopPort,
			MAC:            nil,
			Tag:            tagId,
		}

		bytes, err := json.Marshal(cfd)
		if err != nil {
			return nil, err
		}

		cfl := &global.CircuitSetupLayer{
			Data:         bytes,
			EphPublicKey: MarshalPublicKey(publicKey),
		}

		bytes, err = json.Marshal(cfl)
		if err != nil {
			return nil, err
		}

		message = bytes // set message to encrypted value for next iteration
	}

	return message, nil
}

func DecryptOnion(data []byte, ckt global.Circuit) ([]byte, error) {
	var err error
	for _, e := range ckt {
		data, err = DecryptSymmetric(data, e.SymmetricKey)
		if err != nil {
			return nil, err
		}
	}

	return data, err
}

func DecryptCircuitSetupLayer(data []byte, symmetricKey []byte) (*global.CircuitForwardingData, error) {
	decryptedData, err := DecryptSymmetric(data, symmetricKey)
	if err != nil {
		return nil, err
	}

	cfd := global.ReadCircuitForwardingData(decryptedData)

	return cfd, nil
}
