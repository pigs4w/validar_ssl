package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func scanURL(url string, wg *sync.WaitGroup, ch chan struct{}) {
	defer wg.Done()
	fmt.Printf("Scanning URL: %s\n", url)
	cmd := exec.Command("nmap", "-p", "443", "--script", "ssl-cert", url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Erro ao escanear URL %s: %s\n", url, err)
	}
	<-ch // Liberando o canal para a próxima goroutine
}

func main() {
	fmt.Println("Automatizando o Nmap para URLs...")

	if len(os.Args) != 2 {
		fmt.Println("Uso: go run main.go <arquivo_de_urls>")
		return
	}

	filename := os.Args[1]

	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Erro ao abrir o arquivo %s: %s\n", filename, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var wg sync.WaitGroup
	// Limitando o número de goroutines simultâneas para evitar sobrecarga
	maxGoroutines := 10
	ch := make(chan struct{}, maxGoroutines)

	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		wg.Add(1)
		ch <- struct{}{} // Consumindo o canal para limitar o número de goroutines
		go scanURL(url, &wg, ch)
	}

	wg.Wait()

	fmt.Println("Concluído!")
}
