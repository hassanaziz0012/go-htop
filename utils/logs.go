package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

func WriteLog(msg string) {
	f, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fmt.Fprintf(f, "[%s]%s\n", time.Now(), msg)
}
