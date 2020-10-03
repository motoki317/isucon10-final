package main

import (
	"crypto/elliptic"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

const (
	WebpushVAPIDPrivateKeyPath = "../vapid_private.pem"
	WebpushSubject             = "xsuportal@example.com"
)

func main() {
	pemBytes, err := ioutil.ReadFile(WebpushVAPIDPrivateKeyPath)
	if err != nil {
		return
	}
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return
	}
	priKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return
	}
	priBytes := priKey.D.Bytes()
	pubBytes := elliptic.Marshal(priKey.Curve, priKey.X, priKey.Y)
	pri := base64.RawURLEncoding.EncodeToString(priBytes)
	pub := base64.RawURLEncoding.EncodeToString(pubBytes)
	fmt.Println(pri)
	fmt.Println(pub)
	return
}
