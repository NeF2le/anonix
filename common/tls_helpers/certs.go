package tls_helpers

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)

type KeyPair struct {
	PublicKey  *x509.Certificate
	PrivateKey *rsa.PrivateKey
}

func GenerateAndSignCertificate(root *KeyPair, publicKeyFile, privateKeyFile string, dnsNames []string) error {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return err
	}

	cert := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"anonix"},
			Country:      []string{"RU"},
			Province:     []string{"VologdaOblast"},
			Locality:     []string{"Vologda"},
			CommonName:   "localhost",
		},
		DNSNames:     dnsNames,
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1"), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth,
			x509.ExtKeyUsageServerAuth,
		},
		KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, root.PublicKey, &certPrivKey.PublicKey, root.PrivateKey)
	if err != nil {
		return err
	}

	certPEM := new(bytes.Buffer)
	err = pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err != nil {
		return err
	}

	certPrivKeyPEM := new(bytes.Buffer)
	err = pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})
	if err != nil {
		return err
	}

	if err = os.WriteFile(publicKeyFile, certPEM.Bytes(), 0644); err != nil {
		return err
	}

	if err = os.WriteFile(privateKeyFile, certPrivKeyPEM.Bytes(), 0600); err != nil {
		return err
	}

	return nil
}

func ReadCertificate(publicKeyFile, privateKeyFile string) (*KeyPair, error) {
	cert := new(KeyPair)

	privKey, err := os.ReadFile(privateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %v", err)
	}

	privPemBlock, _ := pem.Decode(privKey)
	parsedPrivKey, err := x509.ParsePKCS1PrivateKey(privPemBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	cert.PrivateKey = parsedPrivKey

	pubKey, err := os.ReadFile(publicKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key: %v", err)
	}

	pubPemBlock, _ := pem.Decode(pubKey)
	parsedPubKey, err := x509.ParseCertificate(pubPemBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	cert.PublicKey = parsedPubKey
	return cert, nil
}

func ReadCertificateAuthority(publicKeyFile, privateKeyFile string) (*KeyPair, error) {
	root := new(KeyPair)

	rootKey, err := os.ReadFile(privateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}

	privPemBlock, _ := pem.Decode(rootKey)

	rootPrivKey, err := x509.ParsePKCS8PrivateKey(privPemBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	root.PrivateKey = rootPrivKey.(*rsa.PrivateKey)

	rootCert, err := os.ReadFile(publicKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key: %w", err)
	}

	publicPemBlock, _ := pem.Decode(rootCert)

	rootPubCrt, err := x509.ParseCertificate(publicPemBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	root.PublicKey = rootPubCrt

	return root, nil
}

func Verification(serverName string, cfg *Config) error {
	ca, err := ReadCertificateAuthority(cfg.RootPublicKey, cfg.RootPrivateKey)
	if err != nil {
		return fmt.Errorf("failed to read root CA certificate: %w", err)
	}

	if cfg.AllowAutoGenerate {
		if err = GenerateAndSignCertificate(ca, cfg.ClientPublicKey, cfg.ClientPrivateKey, cfg.DNSNames); err != nil {
			return fmt.Errorf("failed to generate client certificate: %w", err)
		}
		if err = GenerateAndSignCertificate(ca, cfg.ServerPublicKey, cfg.ServerPrivateKey, cfg.DNSNames); err != nil {
			return fmt.Errorf("failed to generate server certificate: %w", err)
		}
	}

	clientCert, err := ReadCertificate(cfg.ClientPublicKey, cfg.ClientPrivateKey)
	if err != nil {
		return fmt.Errorf("failed to read client certificate: %w", err)
	}

	serverCert, err := ReadCertificate(cfg.ServerPublicKey, cfg.ServerPrivateKey)
	if err != nil {
		return fmt.Errorf("failed to read server certificate: %w", err)
	}

	roots := x509.NewCertPool()
	roots.AddCert(ca.PublicKey)

	opts := x509.VerifyOptions{
		Roots:         roots,
		Intermediates: x509.NewCertPool(),
		DNSName:       serverName,
	}

	if _, err = clientCert.PublicKey.Verify(opts); err != nil {
		return fmt.Errorf("failed to verify client certificate: %w", err)
	}

	if _, err = serverCert.PublicKey.Verify(opts); err != nil {
		return fmt.Errorf("failed to verify server certificate: %w", err)
	}

	return nil
}
