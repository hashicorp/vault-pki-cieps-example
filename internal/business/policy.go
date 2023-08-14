package business

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"fmt"
	"math/big"
	"time"

	"github.com/hashicorp/vault/sdk/helper/certutil"
)

func Evaluate(req *certutil.CIEPSRequest) (*certutil.CIEPSResponse, error) {
	var err error

	csr := req.ParsedCSR
	cert := &x509.Certificate{}

	cert.RawSubjectPublicKeyInfo = csr.RawSubjectPublicKeyInfo
	cert.PublicKeyAlgorithm = csr.PublicKeyAlgorithm
	cert.PublicKey = csr.PublicKey

	cert.RawSubject = csr.RawSubject
	cert.Subject = csr.Subject

	cert.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageDataEncipherment | x509.KeyUsageKeyAgreement
	cert.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}

	cert.NotBefore = time.Now().Add(-30 * time.Second)
	cert.NotAfter = time.Now().Add(10 * 24 * time.Hour)

	cert.BasicConstraintsValid = false
	cert.IsCA = false

	cert.SubjectKeyId, err = certutil.GetSubjectKeyID(cert.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to compute subjectKeyId: %w", err)
	}

	cert.SerialNumber = big.NewInt(2)

	extMsg := "CIEPS Demo Server Certificate"
	extValue, err := asn1.Marshal(extMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal certificate extension: %w", err)
	}

	cert.ExtraExtensions = []pkix.Extension{
		{
			Id:    asn1.ObjectIdentifier{2, 16, 840, 1, 113730, 1, 13},
			Value: extValue,
		},
	}

	parent, parentKey, err := getSelfSignedRoot()
	if err != nil {
		return nil, fmt.Errorf("failed to get self-signed root: %w", err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, parent, cert.PublicKey, parentKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal certificate template: %w", err)
	}

	resp := &certutil.CIEPSResponse{
		Warnings: []string{"result from demo server; no validation occurred"},
		ParsedCertificate: &x509.Certificate{
			Raw: certBytes,
		},
		StoreCert: true,
		IssuerRef: req.VaultRequestKV.IssuerID,
	}

	resp.MarshalCertificate()

	return resp, nil
}

func getSelfSignedRoot() (*x509.Certificate, crypto.Signer, error) {
	// XXX: this self-signed root could be persisted.
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate private key: %w")
	}
	pub := key.Public()

	caTemplate := &x509.Certificate{
		Subject:      pkix.Name{CommonName: "CIEPS Demo Root CA"},
		SerialNumber: big.NewInt(1),
		NotBefore:    time.Now(),
		NotAfter:     time.Now(),
		KeyUsage:     x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature | x509.KeyUsageCRLSign,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, pub, key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal CA certificate: %w", err)
	}

	ca, err := x509.ParseCertificate(caBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal CA certificate: %w", err)
	}

	return ca, key, nil
}
