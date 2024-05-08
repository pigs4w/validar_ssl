package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

func validarCertificado(url string) string {
	if !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   5 * time.Second,
	}

	response, err := client.Get(url)
	if err != nil {
		return fmt.Sprintf("Erro ao tentar validar o certificado: %v", err)
	}
	defer response.Body.Close()

	if len(response.TLS.PeerCertificates) == 0 {
		return "Certificado inválido ou não encontrado."
	}

	cert := response.TLS.PeerCertificates[0]

	if cert != nil {
		return fmt.Sprintf("Certificado válido para: %s", cert.Subject.CommonName)
	} else {
		return "Certificado inválido ou não encontrado."
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Uso: go run main.go <arquivo_de_urls>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Erro ao abrir o arquivo: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := scanner.Text()
		resultado := validarCertificado(url)
		fmt.Printf("%s: %s\n", url, resultado)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Erro ao ler o arquivo: %v\n", err)
		os.Exit(1)
	}
}
