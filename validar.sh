#!/bin/bash

# Caminho para o binário Go
BIN_PATH="/caminho/para/seu/binario/go"

# Função para validar uma URL e salvar o resultado em arquivos separados
validate_url() {
    url="$1"
    result=$("$BIN_PATH" "$url")
    if [[ $result == *"Certificado válido"* ]]; then
        echo "URL: $url - Resultado: $result" >> validas.txt
    else
        echo "URL: $url - Resultado: $result" >> invalidas.txt
    fi
}

while IFS= read -r url; do
    validate_url "$url"
done < "urls.txt"
