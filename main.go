package main

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	action    string
	cfg       string
	ca        string
	caKey     string
	outDir    string
	outPrefix string
)

func init() {
	pwd, _ := os.Getwd()
	flag.StringVar(&action, "action", "", "the action ca or sign")
	flag.StringVar(&ca, "ca", "", "ca pem")
	flag.StringVar(&caKey, "ca-key", "", "ca key pem")
	flag.StringVar(&cfg, "cfg", "", "config json file")
	flag.StringVar(&outDir, "out", pwd, "cert and cert-key output dir, default current dir")
	flag.StringVar(&outPrefix, "prefix", "ca", "cert and cert-key filename prefix")
	flag.Usage = func() {
		fmt.Printf(`Usage: sslutil -action [-ca] [-ca-key] -cfg -out -prefix
Options:
`)
		flag.PrintDefaults()
	}
}

// CreateCACert 创建根证书
func CreateCACert(outdir string, prefix string, cfg string) error {
	baseinfo, err := loadInfo(cfg)
	if err != nil {
		return err
	}
	crtinfo := baseinfo
	crtinfo.IsCA = true
	cert, err := CreateCRT(nil, nil, crtinfo)
	if err != nil {
		return err
	}
	return writeOut(outdir, prefix, cert)
}

// CreateSignCert 创建签名证书
func CreateSignCert(ca string, caKey string, outdir string, prefix string, cfg string) error {
	caCertPem, err := ioutil.ReadFile(ca)
	if err != nil {
		return err
	}
	caKeyPem, err := ioutil.ReadFile(caKey)
	if err != nil {
		return err
	}
	keyBlock, _ := pem.Decode(caKeyPem)
	pemBlock, _ := pem.Decode(caCertPem)
	if keyBlock == nil || pemBlock == nil {
		return errors.New("ca key为空 或者 ca pem为空")
	}
	key, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return err
	}
	pem, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		return err
	}
	baseinfo, err := loadInfo(cfg)
	if err != nil {
		return err
	}
	crtinfo := baseinfo
	crtinfo.IsCA = false
	cert, err := CreateCRT(pem, key, crtinfo)
	if err != nil {
		return err
	}
	return writeOut(outdir, prefix, cert)
}

func writeOut(outdir string, prefix string, cert *Cert) error {
	var err error
	if _, err := os.Stat(outdir); os.IsNotExist(err) {
		os.MkdirAll(outdir, 0755)
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s-%s", outdir, prefix, "cert.pem"), cert.Pem, 0755)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s-%s", outdir, prefix, "key.pem"), cert.Key, 0755)
	if err != nil {
		return err
	}
	return nil
}

func loadInfo(cfg string) (CertInformation, error) {
	info := CertInformation{}
	bts, err := ioutil.ReadFile(cfg)
	if err != nil {
		return info, err
	}
	err = json.Unmarshal(bts, &info)
	if err != nil {
		return info, err
	}
	return info, nil
}

func main() {
	var err error
	flag.Parse()
	switch action {
	case "ca":
		err = CreateCACert(outDir, outPrefix, cfg)
	case "sign":
		err = CreateSignCert(ca, caKey, outDir, outPrefix, cfg)
	default:
		flag.Usage()
	}
	if err != nil {
		fmt.Printf("Something error %v\n", err)
	}
}
