package main

import (
	"bufio"
	"log"
	"net"
	"net/http"

	"crypto/tls"
	"crypto/x509"

	"io"

	"golang.org/x/net/http2"
)

func h2c() *http.Client {

	client := &http.Client{}

	client.Transport = &http2.Transport{
		// So http2.Transport doesn't complain the URL scheme isn't 'https'
		AllowHTTP: true,
		// Pretend we are dialing a TLS endpoint.
		// Note, we ignore the passed tls.Config
		DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
			return net.Dial(network, addr)
		},
	}
	return client
}

func h2s() *http.Client {

	client := &http.Client{}
	certBytes := []byte(`-----BEGIN CERTIFICATE-----
MIIDCTCCAfGgAwIBAgIULlv9pzYIhVN14+JmFbGlxGT+B4owDQYJKoZIhvcNAQEL
BQAwFDESMBAGA1UEAwwJbG9jYWxob3N0MB4XDTIxMDQyMDAzMDI0OVoXDTIxMDUy
MDAzMDI0OVowFDESMBAGA1UEAwwJbG9jYWxob3N0MIIBIjANBgkqhkiG9w0BAQEF
AAOCAQ8AMIIBCgKCAQEAzo3hilzPM2hvLcrnqM6mTt4zFM6wc/0Es1VvSVw9oKnh
C/2t2u1321XKe702mgVxLrJAnujk2UjLTwp5qR69IKmlplRNXNLQgPftAcyy/ode
x5Ilej+3aNF5a6k2Mw+jVCYnYOzEW22JQXDZ/97wgGJUK4RWshO0Pr5dINV8juUv
ndw3A8DZCkihCIBTP03oD9RdIAJZ2cfSrW1c0t8g+rjsnsAyrvRAR8juw+kYqDTX
MAP11BGuehoPl7p20TdXinRV2q+bGYapsV+R2SLBnzt3AwjgcdhAbPHHIcP7h0zT
qTjmMPrdhRm1+kbc2FavOSrSobqBsiXiH1pEDf3djwIDAQABo1MwUTAdBgNVHQ4E
FgQU3lWviuzuiXnpC58S4695ztd8o6AwHwYDVR0jBBgwFoAU3lWviuzuiXnpC58S
4695ztd8o6AwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAg5hc
+VeUS0tTo56NfE46tJBkYaFn9fpaYy0DlsNhMQocB3fRh+5Eoawp4mCkmJrmj78f
p4dgnlBxgGOVnNv4whovDp6UYQz+WYKXH68+f8TYB7fZbH3h1q7+/7dTPXl9eFrZ
f66+NnZK6xnlO7b/5E+WscoGnLXdhBtHgktASqWwgVk+OES9Z481lE2j8Kyt6k4a
ePyTJWYH/0rx+U7NMIWXu1UXJm+EGk4V1zl+gqXp7jHQRCE0r5v2jCNtN4GfWDoH
lPz8F2lzhqvuIEWqLgHaRnoPO7PDtsrxFFYbLcFnwDzbeYzZUf/fPRu3yznBXqKh
mWXRIbdLpgKRysZmBg==
-----END CERTIFICATE-----`)

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(certBytes)

	client.Transport = &http2.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:            caCertPool,
			InsecureSkipVerify: true,
		},
	}
	return client
}

func https() *http.Client {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	return http.DefaultClient
}

func CallH2C(url string, onData func(string)) {

	client := https()

	// Perform the request
	resp, err := client.Get(url)
	if err != nil {
		log.Fatalf("Failed calling: %s", err)
	}
	defer resp.Body.Close()

	bufferedReader := bufio.NewReader(resp.Body)

	buffer := make([]byte, 4*1024)

	var totalBytesReceived int

	// Reads the response
	for {
		len, err := bufferedReader.Read(buffer)
		if len > 0 {
			totalBytesReceived += len
			// log.Println(len, "bytes received")
			// Prints received data
			data := string(buffer[:len])
			onData(data)
			// log.Println(data)
		}

		if err != nil {
			if err == io.EOF {
				// Last chunk received
				// log.Println(err)
			}
			break
		}
	}
	log.Println("Total Bytes Received:", totalBytesReceived)
}

// func main() {
// 	CallH2C("https://localhost:8443", func(data string) {
// 		log.Println("Received: " + data)
// 	})
// }
