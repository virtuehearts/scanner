package key

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"golang.org/x/crypto/ripemd160"
)

// KeyChainFromPriv converts a private key to a KeyChain.
// Optimized for Bitcoin P2PKH by using direct hashing for address generation.
func KeyChainFromPriv(priv []byte) (out KeyChain, err error) {
	_, pub := btcec.PrivKeyFromBytes(priv)

	// Fast path for P2PKH compressed address
	compPub := pub.SerializeCompressed()
	compressedAddr, err := hash160ToAddress(compPub)
	if err != nil {
		return out, err
	}

	// Fast path for P2PKH uncompressed address
	uncompPub := pub.SerializeUncompressed()
	uncompressedAddr, err := hash160ToAddress(uncompPub)
	if err != nil {
		return out, err
	}

	return KeyChain{
		Private:      hex.EncodeToString(priv),
		Compressed:   compressedAddr,
		Uncompressed: uncompressedAddr,
	}, nil
}

// hash160ToAddress performs SHA256 + RIPEMD160 + Base58Check
func hash160ToAddress(pub []byte) (string, error) {
	h256 := sha256.Sum256(pub)
	rp := ripemd160.New()
	rp.Write(h256[:])
	h160 := rp.Sum(nil)

	addr, err := btcutil.NewAddressPubKeyHash(h160, &chaincfg.MainNetParams)
	if err != nil {
		return "", err
	}
	return addr.EncodeAddress(), nil
}

// PrivToWIF converts raw private key bytes to WIF format
func PrivToWIF(key []byte) (out string, err error) {
	priv, _ := btcec.PrivKeyFromBytes(key)
	wif, err := btcutil.NewWIF(priv, &chaincfg.MainNetParams, false)
	if err != nil {
		return out, fmt.Errorf("failed to new wif: %w", err)
	}

	return wif.String(), nil
}
