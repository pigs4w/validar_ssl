package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// verifyCertificate verifica o certificado SSL para o URL fornecido.
func verifyCertificate(url string) (bool, string) {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	// Dividindo o URL para extrair o nome do host
	hostName := strings.TrimPrefix(strings.TrimPrefix(url, "http://"), "https://")
	hostName = strings.Split(hostName, "/")[0] // Remover qualquer caminho

	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 10 * time.Second}, "tcp", hostName+":443", nil)
	if err != nil {
		return false, fmt.Sprintf("%s: Erro ao conectar ou ler o certificado\n", url)
	}
	defer conn.Close()

	cert := conn.ConnectionState().PeerCertificates[0]

	// Formatar datas
	validFrom := cert.NotBefore.Format("2006-01-02")
	validUntil := cert.NotAfter.Format("2006-01-02")

	result := fmt.Sprintf("%s\nAssunto: %s\nEmissor: %s\nVálido De: %s\nVálido Até: %s\n",
		url, cert.Subject.CommonName, cert.Issuer.CommonName, validFrom, validUntil)
	return true, result
}

func checkRedirection(url string) bool {
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode >= 300 && resp.StatusCode < 400
}

func processUrl(url string, validFile, invalidFile *os.File, wg *sync.WaitGroup, mutex *sync.Mutex) {
	defer wg.Done()
	valid, result := verifyCertificate(url)
	mutex.Lock()
	if valid {
		validFile.WriteString(result)
	} else {
		if strings.HasPrefix(url, "https://") {
			if checkRedirection(url) {
				mutex.Unlock()
				validFile.WriteString(url + " (Redirecionado)\n")
				return
			}
		}
		invalidFile.WriteString(result)
	}
	mutex.Unlock()
}

func main() {
	fmt.Println("Analisando URLs...")

	urlsFile, err := os.Open("urls.txt")
	if err != nil {
		fmt.Println("Erro ao abrir arquivo de URL:", err)
		return
	}
	defer urlsFile.Close()

	validFile, err := os.Create("validadas.txt")
	if err != nil {
		fmt.Println("Erro ao criar arquivo válido:", err)
		return
	}
	defer validFile.Close()

	invalidFile, err := os.Create("invalidadas.txt")
	if err != nil {
		fmt.Println("Erro ao criar arquivo inválido:", err)
		return
	}
	defer invalidFile.Close()

	scanner := bufio.NewScanner(urlsFile)
	var wg sync.WaitGroup
	var mutex sync.Mutex

	for scanner.Scan() {
		wg.Add(1)
		go processUrl(scanner.Text(), validFile, invalidFile, &wg, &mutex)
	}

	wg.Wait()

	fmt.Println("Concluído!")
}
