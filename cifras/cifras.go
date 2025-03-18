package cifras

import (
	"fmt"
	"strings"
)

type operacion uint8

const (
	sumar operacion = iota
	restar
	multiplicar
	dividir
)

var operaciones = []operacion{sumar, restar, multiplicar, dividir}

func (op operacion) String() string {
	return [...]string{"+", "-", "x", "/"}[op]
}

type nodo struct {
	operando1, operando2 uint32
	calculo              operacion
	resultado            uint32
	operandosPorUsar     []uint32
	padre                *nodo
	hijos                []*nodo
}

func (n *nodo) calcularDescendientes() {
	for i := 0; i < len(n.operandosPorUsar)-1; i++ {
		for j := i + 1; j < len(n.operandosPorUsar); j++ {
			for _, c := range operaciones {
				hijo := nodo{calculo: c, padre: n}
				if n.operandosPorUsar[i] >= n.operandosPorUsar[j] {
					hijo.operando1 = n.operandosPorUsar[i]
					hijo.operando2 = n.operandosPorUsar[j]
				} else {
					hijo.operando1 = n.operandosPorUsar[j]
					hijo.operando2 = n.operandosPorUsar[i]
				}
				resultadoValido := true
				switch c {
				case sumar:
					hijo.resultado = hijo.operando1 + hijo.operando2
				case restar:
					hijo.resultado = hijo.operando1 - hijo.operando2
				case multiplicar:
					hijo.resultado = hijo.operando1 * hijo.operando2
				case dividir:
					if hijo.operando2 != 0 && hijo.operando1%hijo.operando2 == 0 {
						hijo.resultado = hijo.operando1 / hijo.operando2
					} else {
						resultadoValido = false
					}
				}
				if resultadoValido {
					hijo.operandosPorUsar = append(hijo.operandosPorUsar, hijo.resultado)
					hijo.operandosPorUsar = append(hijo.operandosPorUsar, n.operandosPorUsar[:i]...)
					hijo.operandosPorUsar = append(hijo.operandosPorUsar, n.operandosPorUsar[i+1:j]...)
					hijo.operandosPorUsar = append(hijo.operandosPorUsar, n.operandosPorUsar[j+1:]...)
					n.hijos = append(n.hijos, &hijo)
					hijo.calcularDescendientes()
				}
			}
		}
	}
}

func (n *nodo) encontrarAproximacion(objetivo uint32, aproximacion *uint32, candidatos *[]*nodo) {
	var diferencia uint32
	if objetivo > n.resultado {
		diferencia = objetivo - n.resultado
	} else {
		diferencia = n.resultado - objetivo
	}
	if diferencia < *aproximacion {
		*aproximacion = diferencia
		*candidatos = append((*candidatos)[:0], n)
	} else if diferencia == *aproximacion {
		*candidatos = append(*candidatos, n)
	}
	for _, hijo := range n.hijos {
		hijo.encontrarAproximacion(objetivo, aproximacion, candidatos)
	}
}

func (n nodo) String() string {
	if n.padre == nil {
		return ""
	} else if n.padre.padre == nil {
		return fmt.Sprintf("(%d %s %d)", n.operando1, n.calculo, n.operando2)
	} else if n.operando1 == n.padre.resultado {
		return fmt.Sprintf("(%s %s %d)", *n.padre, n.calculo, n.operando2)
	} else if n.operando2 == n.padre.resultado {
		return fmt.Sprintf("(%d %s %s)", n.operando1, n.calculo, *n.padre)
	} else {
		return fmt.Sprintf("(%d %s %d)", n.operando1, n.calculo, n.operando2)
	}
}

func Resuelve(operandos []uint32, resultado uint32) string {
	var raiz = nodo{operandosPorUsar: operandos}
	raiz.calcularDescendientes()
	var mejorAproximacion uint32 = 999
	var nodosCandidatos []*nodo
	raiz.encontrarAproximacion(resultado, &mejorAproximacion, &nodosCandidatos)
	var escritorOperaciones strings.Builder
	if mejorAproximacion == 0 {
		escritorOperaciones.WriteString("Encontrado exacto:\n")
	} else {
		escritorOperaciones.WriteString(fmt.Sprintf("Me quedo a %d:\n", mejorAproximacion))
	}
	nodosEscritos := make([]string, len(nodosCandidatos))
	for i := range nodosCandidatos {
		nodosEscritos[i] = nodosCandidatos[i].String()
	}
	escritorOperaciones.WriteString(strings.Join(nodosEscritos, "\n"))
	return escritorOperaciones.String()
}
