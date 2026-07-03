package crypto

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
)

func SignTransactionHash(hashBase64 string, privateKeyBase64 string) (string, error) {
	hashBytes, err := base64.StdEncoding.DecodeString(hashBase64)
	if err != nil {
		return "", fmt.Errorf("failed to decode hash: %w", err)
	}

	privateKeyBytes, err := base64.StdEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		return "", fmt.Errorf("failed to decode private key: %w", err)
	}

	if len(privateKeyBytes) != ed25519.PrivateKeySize {
		return "", fmt.Errorf("invalid private key size: expected %d, got %d", ed25519.PrivateKeySize, len(privateKeyBytes))
	}

	privateKey := ed25519.PrivateKey(privateKeyBytes)
	signature := ed25519.Sign(privateKey, hashBytes)

	return base64.StdEncoding.EncodeToString(signature), nil
}

func VerifySignedTxHash(hashBase64 string, publicKeyBase64 string, signatureBase64 string) bool {
	hashBytes, err := base64.StdEncoding.DecodeString(hashBase64)
	if err != nil {
		return false
	}

	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return false
	}

	signatureBytes, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		return false
	}

	if len(publicKeyBytes) != ed25519.PublicKeySize {
		return false
	}

	publicKey := ed25519.PublicKey(publicKeyBytes)
	return ed25519.Verify(publicKey, hashBytes, signatureBytes)
}
