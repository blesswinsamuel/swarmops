package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"golang.org/x/crypto/ssh"
)

const (
	publicKeyFile  = "key.pub"
	privateKeyFile = "key.pem"
)

func generateRsaKeyPairs(keysDir string) error {
	os.MkdirAll(keysDir, os.ModePerm)
	reader := rand.Reader
	bitSize := 2048

	// https://gist.github.com/sdorra/1c95de8cb80da31610d2ad767cd6f251
	key, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {
		return err
	}

	err = savePublicSSHKey(path.Join(keysDir, publicKeyFile), key.PublicKey)
	if err != nil {
		return err
	}
	err = savePrivatePEMKey(path.Join(keysDir, privateKeyFile), key)
	if err != nil {
		return err
	}
	return nil
}

func getPublicKey() (string, error) {
	publicKey, err := ioutil.ReadFile(path.Join(*keysDir, publicKeyFile))
	if err != nil {
		return "", err
	}
	return string(publicKey), nil
}

func savePrivatePEMKey(fileName string, key *rsa.PrivateKey) error {
	outFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer outFile.Close()

	var privateKey = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(outFile, privateKey)
	if err != nil {
		return err
	}
	return nil
}

func savePublicSSHKey(fileName string, pubkey rsa.PublicKey) error {
	pub, err := ssh.NewPublicKey(&pubkey)
	if err != nil {
		return fmt.Errorf("ssh.NewPublicKey: %w", err)
	}
	err = ioutil.WriteFile(fileName, ssh.MarshalAuthorizedKey(pub), 0655)
	if err != nil {
		return fmt.Errorf("ioutil.WriteFile: %w", err)
	}
	return nil
}
