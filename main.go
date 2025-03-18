package main

import (
	"fmt"

	"github.com/rodrigargar/cifrasyletras/cifras"
)

func main() {
	var operandos [6]uint32
	var resultado uint32

	fmt.Println("Introduce los seis números:")
	for i := range operandos {
		fmt.Printf("%d: ", i+1)
		fmt.Scanf("%d", &operandos[i])
	}
	fmt.Println("Resultado a conseguir:")
	fmt.Scanf("%d", &resultado)

	fmt.Println(cifras.Resuelve(operandos[:], resultado))
}
