// Copyright 2021 The Casdoor Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package object

import (
	"bytes"
	"compress/flate"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"encoding/xml"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"

	"github.com/RobotsAndPencils/go-saml"
	"github.com/beevik/etree"
	"github.com/casdoor/casdoor/conf"
	saml2 "github.com/russellhaering/gosaml2"
	dsig "github.com/russellhaering/goxmldsig"
)

func ParseSamlResponse(samlResponse string, providerType string) (string, error) {
	samlResponse, _ = url.QueryUnescape(samlResponse)
	sp, err := buildSp(&Provider{Type: providerType}, samlResponse)
	if err != nil {
		return "", err
	}
	assertionInfo, err := sp.RetrieveAssertionInfo(samlResponse)
	if err != nil {
		panic(err)
	}
	return assertionInfo.NameID, nil
}

func GenerateSamlLoginUrl(id, relayState string) (string, string, error) {
	provider := GetProvider(id)
	if provider.Category != "SAML" {
		return "", "", fmt.Errorf("Provider %s's category is not SAML", provider.Name)
	}
	sp, err := buildSp(provider, "")
	if err != nil {
		return "", "", err
	}
	auth := ""
	method := ""
	if provider.EnableSignAuthnRequest {
		post, err := sp.BuildAuthBodyPost(relayState)
		if err != nil {
			return "", "", err
		}
		auth = string(post[:])
		method = "POST"
	} else {
		auth, err = sp.BuildAuthURL(relayState)
		if err != nil {
			return "", "", err
		}
		method = "GET"
	}
	return auth, method, nil
}

func buildSp(provider *Provider, samlResponse string) (*saml2.SAMLServiceProvider, error) {
	certStore := dsig.MemoryX509CertificateStore{
		Roots: []*x509.Certificate{},
	}
	origin := conf.GetConfigString("origin")
	certEncodedData := ""
	if samlResponse != "" {
		certEncodedData = parseSamlResponse(samlResponse, provider.Type)
	} else if provider.IdP != "" {
		certEncodedData = provider.IdP
	}
	certData, err := base64.StdEncoding.DecodeString(certEncodedData)
	if err != nil {
		return nil, err
	}
	idpCert, err := x509.ParseCertificate(certData)
	if err != nil {
		return nil, err
	}
	certStore.Roots = append(certStore.Roots, idpCert)
	sp := &saml2.SAMLServiceProvider{
		ServiceProviderIssuer:       fmt.Sprintf("%s/api/acs", origin),
		AssertionConsumerServiceURL: fmt.Sprintf("%s/api/acs", origin),
		IDPCertificateStore:         &certStore,
		SignAuthnRequests:           false,
		SPKeyStore:                  dsig.RandomKeyStoreForTest(),
	}
	if provider.Endpoint != "" {
		sp.IdentityProviderSSOURL = provider.Endpoint
		sp.IdentityProviderIssuer = provider.IssuerUrl
	}
	if provider.EnableSignAuthnRequest {
		sp.SignAuthnRequests = true
		sp.SPKeyStore = buildSpKeyStore()
	}
	return sp, nil
}

func parseSamlResponse(samlResponse string, providerType string) string {
	de, err := base64.StdEncoding.DecodeString(samlResponse)
	if err != nil {
		panic(err)
	}
	deStr := strings.Replace(string(de), "\n", "", -1)
	tagMap := map[string]string{
		"Aliyun IDaaS": "ds",
		"Keycloak":     "dsig",
	}
	tag := tagMap[providerType]
	expression := fmt.Sprintf("<%s:X509Certificate>([\\s\\S]*?)</%s:X509Certificate>", tag, tag)
	res := regexp.MustCompile(expression).FindStringSubmatch(deStr)
	return res[1]
}

func buildSpKeyStore() dsig.X509KeyStore {
	keyPair, err := tls.LoadX509KeyPair("object/token_jwt_key.pem", "object/token_jwt_key.key")
	if err != nil {
		panic(err)
	}
	return &dsig.TLSCertKeyStore{
		PrivateKey:  keyPair.PrivateKey,
		Certificate: keyPair.Certificate,
	}
}
func GetSamlResponse(application *Application, user *User, samlRequest string, host string) (string, string, error) {
	//decode samlRequest
	defated, err := base64.StdEncoding.DecodeString(samlRequest)
	if err != nil {
		return "", "", fmt.Errorf("err: %s", err.Error())
	}
	var buffer bytes.Buffer
	rdr := flate.NewReader(bytes.NewReader(defated))
	io.Copy(&buffer, rdr)
	var authnRequest saml.AuthnRequest
	err = xml.Unmarshal(buffer.Bytes(), &authnRequest)
	if err != nil {
		return "", "", fmt.Errorf("err: %s", err.Error())
	}
	//verify samlRequest
	if valid := CheckRedirectUriValid(application, authnRequest.Issuer.Url); !valid {
		return "", "", fmt.Errorf("err: invalid issuer url")
	}

	//get publickey string
	cert := getCertByApplication(application)
	block, _ := pem.Decode([]byte(cert.PublicKey))
	publicKey := base64.StdEncoding.EncodeToString(block.Bytes)

	_, originBackend := getOriginFromHost(host)

	//build signedResponse
	samlResponse, _ := NewSamlResponse(user, originBackend, publicKey, authnRequest.AssertionConsumerServiceURL, authnRequest.Issuer.Url, application.RedirectUris)
	randomKeyStore := &X509Key{
		PrivateKey:      cert.PrivateKey,
		X509Certificate: publicKey,
	}
	ctx := dsig.NewDefaultSigningContext(randomKeyStore)
	ctx.Hash = crypto.SHA1
	signedXML, err := ctx.SignEnveloped(samlResponse)
	if err != nil {
		return "", "", fmt.Errorf("err: %s", err.Error())
	}

	doc := etree.NewDocument()
	doc.SetRoot(signedXML)
	xmlStr, err := doc.WriteToString()
	if err != nil {
		return "", "", fmt.Errorf("err: %s", err.Error())
	}
	res := base64.StdEncoding.EncodeToString([]byte(xmlStr))
	return res, authnRequest.AssertionConsumerServiceURL, nil
}
