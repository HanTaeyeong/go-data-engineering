package utils

import (
	"bytes"
	"io"
	"log"
	"os"
	"strings"
)

func Parse(r io.Reader) string {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)

	if err != nil {
		log.Fatal(err)
	}

	return buf.String()
}

func ReadAndSetEnvFile(filename string) {
	file, err := os.Open(filename)

	if err != nil {
		return
	}
	defer file.Close()
	fileString := Parse(file)

	rows := strings.Split(fileString, "\n")

	for _, row := range rows {
		keyVal := strings.Split(row, "=")
		key := keyVal[0]
		value := keyVal[1]

		os.Setenv(key, value)
	}
}
