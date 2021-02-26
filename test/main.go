package main

import (
	"fmt"

	"github.com/AlexandrGurkin/http_template/internal/ver"
)

func main() {
	fmt.Print("version:", ver.GetVersion())
}
