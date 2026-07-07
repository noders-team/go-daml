package crypto

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
)

const (
	CantonHashPurposePublicKeyFingerprint  = 12
	CantonHashPurposePreparedTransaction   = 48
	CantonHashPurposeMultiTopologyTxHashes = 55
)

func ComputeSHA256CantonHash(purpose int, data []byte) ([]byte, error) {
	h := sha256.New()

	purposeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(purposeBytes, uint32(purpose))

	if _, err := h.Write(purposeBytes); err != nil {
		return nil, fmt.Errorf("failed to write purpose to hash: %w", err)
	}

	if _, err := h.Write(data); err != nil {
		return nil, fmt.Errorf("failed to write data to hash: %w", err)
	}

	hash := h.Sum(nil)

	multiprefix := []byte{0x12, 0x20}
	fullHash := append(multiprefix, hash...)

	return fullHash, nil
}

func ComputeMultiHashForTopology(hashes [][]byte) ([]byte, error) {
	h := sha256.New()
	for _, hash := range hashes {
		h.Write(hash)
	}
	return h.Sum(nil), nil
}

func HashPreparedTransaction(preparedTransactionBase64 string) (string, error) {
	preparedBytes, err := base64.StdEncoding.DecodeString(preparedTransactionBase64)
	if err != nil {
		return "", fmt.Errorf("failed to decode prepared transaction: %w", err)
	}

	hash, err := ComputeSHA256CantonHash(CantonHashPurposePreparedTransaction, preparedBytes)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(hash), nil
}
