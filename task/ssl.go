package task

import (
	"ampctl/config"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"

	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"
)

type GenerateRootCaTask struct {
	Config *config.Config
}

func (t *GenerateRootCaTask) Run() error {
	_, err := os.OpenFile(t.Config.Apache.SslCertificateFile, os.O_RDONLY, 0600)
	if err == nil {
		fmt.Println("Certificate already exits")
		return nil
	}

	subject := pkix.Name{
		CommonName:         t.Config.Apache.SslCertificateCn,
		Country:            []string{t.Config.Apache.SslCertificateCountry},
		Organization:       []string{t.Config.Apache.SslCertificateOrganization},
		OrganizationalUnit: []string{t.Config.Apache.SslCertificateOrganizationUnit},
		Locality:           []string{t.Config.Apache.SslCertificateLocality},
		Province:           []string{t.Config.Apache.SslCertificateProvince},
	}

	// Key parameters
	keyBits := 4096

	// Validity
	notBefore := time.Now().Add(-5 * time.Minute)        // backdate slightly
	notAfter := notBefore.Add(10 * 365 * 24 * time.Hour) // ~10 years

	// Serial number
	serialLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialLimit)
	if err != nil {
		return fmt.Errorf("serial: %w", err)
	}

	// Create CA template
	tpl := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      subject,
		NotBefore:    notBefore,
		NotAfter:     notAfter,

		// Mark as CA
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		// No SANs for a root CA; leave DNS/IPs empty.
		// PathLen=0 would restrict to signing end-entity certs only; omit or set >0 to allow intermediates.
		MaxPathLenZero: false, // allow intermediates if you later want them
	}

	// Generate private key
	privKey, err := rsa.GenerateKey(rand.Reader, keyBits)
	if err != nil {
		return fmt.Errorf("key gen: %w", err)
	}

	// Self-sign (parent = template, signer = our key)
	derBytes, err := x509.CreateCertificate(rand.Reader, tpl, tpl, &privKey.PublicKey, privKey)
	if err != nil {
		return fmt.Errorf("create cert: %w", err)
	}

	// Write certificate (PEM)
	if err := writePem(t.Config.Apache.SslCertificateFile, "CERTIFICATE", derBytes, 0644); err != nil {
		return fmt.Errorf("write cert: %w", err)
	}

	// Write private key (PKCS#1 PEM)
	keyFile, err := os.OpenFile(t.Config.Apache.SslCertificateKeyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("open key: %w", err)
	}
	defer keyFile.Close()
	if err := pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privKey)}); err != nil {
		return fmt.Errorf("encode key: %w", err)
	}

	log.Println("Wrote CA.pem and CA.key")

	return nil
}

func writePem(path, blockType string, der []byte, mode os.FileMode) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer f.Close()
	return pem.Encode(f, &pem.Block{Type: blockType, Bytes: der})
}

type GenerateHostsCaTask struct {
	Config *config.Config
}

func (t *GenerateHostsCaTask) Run() error {
	// Ensure the ssl directory exists before generating host certificates
	if err := os.MkdirAll("/opt/homebrew/etc/httpd/ssl", 0755); err != nil {
		return err
	}
	for _, host := range t.Config.Hosts {
		if host.Ssl {
			if err := t.createHostCertificate(host.Host); err != nil {
				return err
			}
		}
	}

	// default certificate
	if err := t.createHostCertificate("localhost"); err != nil {
		return err
	}

	return nil
}

func (t *GenerateHostsCaTask) createHostCertificate(hostname string) error {
	hostCertFilePath := fmt.Sprintf("/opt/homebrew/etc/httpd/ssl/%s.pem", hostname)
	_, err := os.OpenFile(hostCertFilePath, os.O_RDONLY, 0600)
	if err == nil {
		return nil
	}

	fmt.Printf("Creating host certificate for %s\n", hostname)

	// Load root CA certificate
	caCertPEM, err := os.ReadFile(t.Config.Apache.SslCertificateFile)
	if err != nil {
		return fmt.Errorf("read CA cert: %w", err)
	}
	caKeyPEM, err := os.ReadFile(t.Config.Apache.SslCertificateKeyFile)
	if err != nil {
		return fmt.Errorf("read CA key: %w", err)
	}

	block, _ := pem.Decode(caCertPEM)
	if block == nil || block.Type != "CERTIFICATE" {
		return fmt.Errorf("failed to decode CA cert PEM")
	}
	caCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("parse CA cert: %w", err)
	}

	block, _ = pem.Decode(caKeyPEM)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return fmt.Errorf("failed to decode CA key PEM")
	}
	caKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("parse CA key: %w", err)
	}

	// Generate host private key
	hostKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("generate host key: %w", err)
	}

	// Certificate template for host
	serial, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	hostTemplate := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:         hostname,
			Country:            []string{t.Config.Apache.SslCertificateCountry},
			Organization:       []string{t.Config.Apache.SslCertificateOrganization},
			OrganizationalUnit: []string{t.Config.Apache.SslCertificateOrganizationUnit},
			Locality:           []string{t.Config.Apache.SslCertificateLocality},
			Province:           []string{t.Config.Apache.SslCertificateProvince},
		},
		NotBefore: time.Now().Add(-time.Hour),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour),

		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,

		DNSNames: []string{hostname},
	}

	// Sign host cert with CA
	derBytes, err := x509.CreateCertificate(rand.Reader, hostTemplate, caCert, &hostKey.PublicKey, caKey)
	if err != nil {
		return fmt.Errorf("create host cert: %w", err)
	}

	// Write host cert
	hostCertFile, _ := os.Create(hostCertFilePath)
	defer hostCertFile.Close()
	pem.Encode(hostCertFile, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	// append CA cert so clients see full chain
	pem.Encode(hostCertFile, &pem.Block{Type: "CERTIFICATE", Bytes: caCert.Raw})

	// Write host private key
	hostKeyFilePath := fmt.Sprintf("/opt/homebrew/etc/httpd/ssl/%s.key", hostname)
	hostKeyFile, _ := os.Create(hostKeyFilePath)
	defer hostKeyFile.Close()
	pem.Encode(hostKeyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(hostKey)})

	fmt.Println("Wrote host.pem (host cert + CA) and host.key")
	return nil
}
