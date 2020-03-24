package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	privkey, pubkey := GenerateKeyPair(2048)
	privBytes := PrivateKeyToBytes(privkey)
	pubBytes := PublicKeyToBytes(pubkey)

	WritePrivatePEM("private_key.pem", privBytes)
	WritePublicPEM("public_key.pem", pubBytes)
}

func WritePrivatePEM(filePath string, privBytes *pem.Block) {
	pemPrivateFile, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = pem.Encode(pemPrivateFile, privBytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	pemPrivateFile.Close()
}

func WritePublicPEM(filePath string, pubBytes *pem.Block) {
	pemPublicFile, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = pem.Encode(pemPublicFile, pubBytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	pemPublicFile.Close()
}

// GenerateKeyPair generates a new key pair
func GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		log.Error(err)
	}
	return privkey, &privkey.PublicKey
}

// PrivateKeyToBytes private key to bytes
func PrivateKeyToBytes(priv *rsa.PrivateKey) *pem.Block {
	privBytes :=  &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		}

	return privBytes
}

// PublicKeyToBytes public key to bytes
func PublicKeyToBytes(pub *rsa.PublicKey) *pem.Block {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		log.Error(err)
	}

	pubBytes := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	}

	return pubBytes
}

// BytesToPrivateKey bytes to private key
func BytesToPrivateKey(priv []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			log.Error(err)
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		log.Error(err)
	}
	return key
}

// BytesToPublicKey bytes to public key
func BytesToPublicKey(pub []byte) *rsa.PublicKey {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			log.Error(err)
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		log.Error(err)
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		log.Error("not ok")
	}
	return key
}

// EncryptWithPublicKey encrypts data with public key
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) []byte {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		log.Error(err)
	}
	return ciphertext
}

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) []byte {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if err != nil {
		log.Error(err)
	}
	return plaintext
}