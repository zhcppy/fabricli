/*
@Time 2019-09-06 18:30
@Author ZH

*/
package actions

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/pkg/errors"
)

func ParsePubKeyFromCert(cert []byte) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode(cert)

	if block.Type == "NEW CERTIFICATE REQUEST" || block.Type == "CERTIFICATE REQUEST" {
		csrReq, err := x509.ParseCertificateRequest(block.Bytes)
		if err != nil {
			return nil, err
		}
		lowLevelKey, ok := csrReq.PublicKey.(*ecdsa.PublicKey)
		if !ok {
			return nil, errors.New("invalid raw material. Expected *ecdsa.PublicKey")
		}
		return lowLevelKey, nil
	} else if block.Type == "CERTIFICATE" {
		x509Cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}
		lowLevelKey, ok := x509Cert.PublicKey.(*ecdsa.PublicKey)
		if !ok {
			return nil, errors.New("invalid raw material. Expected *ecdsa.PublicKey")
		}
		return lowLevelKey, nil
	}
	return nil, errors.New(block.Type + " not support")
}
