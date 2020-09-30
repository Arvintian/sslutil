package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	rd "math/rand"
	"net"
	"net/mail"
	"net/url"
	"time"
)

type CertInformation struct {
	Country            []string `json:"country"`
	Organization       []string `json:"organization"`
	OrganizationalUnit []string `json:"organizationalUnit"`
	Province           []string `json:"province"`
	Locality           []string `json:"locality"`
	CrtName            string
	KeyName            string
	IsCA               bool
	Hosts              []string `json:"hosts"`
	Years              int      `json:"years"`
}

type Cert struct {
	Pem []byte
	Key []byte
}

func CreateCRT(RootCa *x509.Certificate, RootKey *rsa.PrivateKey, info CertInformation) (*Cert, error) {
	Certificate := NewCertificate(info)
	Key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	var buf []byte
	if RootCa == nil || RootKey == nil {
		//创建自签名证书
		buf, err = x509.CreateCertificate(rand.Reader, Certificate, Certificate, &Key.PublicKey, Key)
	} else {
		//使用根证书签名
		buf, err = x509.CreateCertificate(rand.Reader, Certificate, RootCa, &Key.PublicKey, RootKey)
	}
	if err != nil {
		return nil, err
	}
	var cer *pem.Block = &pem.Block{Bytes: buf, Type: "CERTIFICATE"}
	keybuf := x509.MarshalPKCS1PrivateKey(Key)
	var key *pem.Block = &pem.Block{Bytes: keybuf, Type: "PRIVATE KEY"}

	c := &Cert{
		Pem: pem.EncodeToMemory(cer),
		Key: pem.EncodeToMemory(key),
	}
	return c, nil
}

func NewCertificate(info CertInformation) *x509.Certificate {
	var tpl = x509.Certificate{
		SerialNumber: big.NewInt(rd.Int63()),
		Subject: pkix.Name{
			Country:            info.Country,
			Organization:       info.Organization,
			OrganizationalUnit: info.OrganizationalUnit,
			Province:           info.Province,
			Locality:           info.Locality,
		},
		NotBefore:             time.Now(),                                                                 //证书的开始时间
		NotAfter:              time.Now().AddDate(info.Years, 0, 0),                                       //证书的结束时间
		BasicConstraintsValid: true,                                                                       //基本的有效性约束
		IsCA:                  info.IsCA,                                                                  //是否是根证书
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}, //证书用途
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}
	for i := range info.Hosts {
		if ip := net.ParseIP(info.Hosts[i]); ip != nil {
			tpl.IPAddresses = append(tpl.IPAddresses, ip)
		} else if email, err := mail.ParseAddress(info.Hosts[i]); err == nil && email != nil {
			tpl.EmailAddresses = append(tpl.EmailAddresses, email.Address)
		} else if uri, err := url.ParseRequestURI(info.Hosts[i]); err == nil && uri != nil {
			tpl.URIs = append(tpl.URIs, uri)
		} else {
			tpl.DNSNames = append(tpl.DNSNames, info.Hosts[i])
		}
	}

	return &tpl
}
