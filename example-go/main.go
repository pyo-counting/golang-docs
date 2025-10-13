package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	for i, v := range os.Environ() {
		str := strings.SplitN(v, "=", 2)
		fmt.Printf("%d\tkey: %s\tvalue: %s\n", i, str[0], str[1])
	}
}
