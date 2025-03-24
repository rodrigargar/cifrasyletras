package cifras

import (
	"fmt"
	"strings"
)

type operacion uint8

const (
	nop operacion = iota
	sumar
	restar
	multiplicar
	dividir
)

var operaciones = []operacion{nop, sumar, restar, multiplicar, dividir}

func (op operacion) String() string {
	return [...]string{"nop", "+", "-", "x", "/"}[op]
}

type nodo struct {
	operando1, operando2 *nodo
	calculo              operacion
	valor                uint32
	padre                *nodo
	hijos                []*nodo
	restantes            []*nodo
}

func (n *nodo) calcularDescendientes() {
	progenitores := append(n.restantes, n)
	for i := 0; i < len(progenitores)-1; i++ {
		for j := i + 1; j < len(progenitores); j++ {
			for _, c := range operaciones[sumar:] {
				hijo := nodo{calculo: c, padre: n}
				if progenitores[i].valor >= progenitores[j].valor {
					hijo.operando1 = progenitores[i]
					hijo.operando2 = progenitores[j]
				} else {
					hijo.operando1 = progenitores[j]
					hijo.operando2 = progenitores[i]
				}
				resultadoValido := true
				switch c {
				case sumar:
					hijo.valor = hijo.operando1.valor + hijo.operando2.valor
				case restar:
					hijo.valor = hijo.operando1.valor - hijo.operando2.valor
				case multiplicar:
					hijo.valor = hijo.operando1.valor * hijo.operando2.valor
				case dividir:
					if (hijo.operando2.valor != 0) && (hijo.operando1.valor%hijo.operando2.valor == 0) {
						hijo.valor = hijo.operando1.valor / hijo.operando2.valor
					} else {
						resultadoValido = false
					}
				default:
					resultadoValido = false
				}
				if resultadoValido {
					hijo.restantes = append(hijo.restantes, progenitores[:i]...)
					hijo.restantes = append(hijo.restantes, progenitores[i+1:j]...)
					hijo.restantes = append(hijo.restantes, progenitores[j+1:]...)
					n.hijos = append(n.hijos, &hijo)
					hijo.calcularDescendientes()
				}
			}
		}
	}
}

func (n *nodo) encontrarAproximacion(objetivo uint32, aproximacion *uint32, candidatos *[]*nodo) {
	var diferencia uint32
	if objetivo > n.valor {
		diferencia = objetivo - n.valor
	} else {
		diferencia = n.valor - objetivo
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
	if n.calculo == nop {
		return fmt.Sprint(n.valor)
	} else {
		return fmt.Sprintf("(%s %s %s)", n.operando1, n.calculo, n.operando2)
	}
}

func Resuelve(operandos []uint32, resultado uint32) string {
	nodosIniciales := make([]*nodo, len(operandos))
	for i := range nodosIniciales {
		nodosIniciales[i] = &nodo{calculo: nop, valor: operandos[i]}
	}
	for j := range nodosIniciales {
		nodosIniciales[j].restantes = append(nodosIniciales[j].restantes, nodosIniciales[:j]...)
		nodosIniciales[j].restantes = append(nodosIniciales[j].restantes, nodosIniciales[j+1:]...)
	}
	for _, n := range nodosIniciales {
		n.calcularDescendientes()
	}
	var mejorAproximacion uint32 = 999
	var nodosCandidatos []*nodo
	for _, n := range nodosIniciales {
		n.encontrarAproximacion(resultado, &mejorAproximacion, &nodosCandidatos)
	}
	var escritorOperaciones strings.Builder
	if mejorAproximacion == 0 {
		escritorOperaciones.WriteString("Encontrado exacto:\n")
	} else {
		escritorOperaciones.WriteString(fmt.Sprintf("Me quedo a %d:\n", mejorAproximacion))
	}
	nodosEscritos := map[string]struct{}{}
	for _, candidato := range nodosCandidatos {
		nodosEscritos[candidato.String()] = struct{}{}
	}
	for operacion := range nodosEscritos {
		escritorOperaciones.WriteString(operacion + "\n")
	}
	return escritorOperaciones.String()
}
