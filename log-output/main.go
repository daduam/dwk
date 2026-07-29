package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func main() {
	randomStr := uuid.New()

	for {
		fmt.Printf("%s: %s\n", time.Now().UTC().Format(time.RFC3339Nano), randomStr)
		time.Sleep(5 * time.Second)
	}
}
