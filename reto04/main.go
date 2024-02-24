package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"unicode"

	"github.com/globalsign/mgo/bson"
	"golang.org/x/net/websocket"
)

/*
a) ¿Cuántos artistas se recomiendan en la página?
b) ¿Cuál es el número de artistas que empiezan por la letra N?
c) El identificador de cada artista es un código alfanumérico. Si cuentas todas las letras de manera independiente que aparecen en cada uno, ¿cuál es el total de caracteres alfabéticos encontrontrados?
*/

func main() {
	var websocketUrl string
	flag.StringVar(&websocketUrl, "url", "wss://ninjascrapers-production.up.railway.app/ws", "url base de la api para el reto")
	flag.Parse()

	wsConexion, wsConexionError := websocket.Dial(websocketUrl, "", websocketUrl)
	if wsConexionError != nil {
		log.Fatalln("error al establecer la conexión al websocket:", wsConexionError)
	}
	defer wsConexion.Close()

	mensajeInicio := MensajeWS{
		Identificador: "artistas",
		Espera:        3,
		Fin:           false,
	}

	mensajeInicioBson, mensajeInicioBsonError := bson.Marshal(mensajeInicio)
	if mensajeInicioBsonError != nil {
		log.Fatalln("error al serializar el mensaje BSON:", mensajeInicioBsonError)

	}

	_, mensajeInicioEnvioError := wsConexion.Write(mensajeInicioBson)
	if mensajeInicioEnvioError != nil {
		log.Fatalln("error al enviar mensaje BSON:", mensajeInicioEnvioError)
	}

	var artistasRecomendados int
	var artistasN int
	var letrasIdentificador int

	for {
		var respuestaRecibida []byte
		respuestaRecibidaError := websocket.Message.Receive(wsConexion, &respuestaRecibida)
		if respuestaRecibidaError != nil {
			log.Fatalln("error al recibir mensaje:", respuestaRecibidaError)
		}

		var respuestaWs RespuestaWs
		respuestaDeserializacionError := bson.Unmarshal(respuestaRecibida, &respuestaWs)
		if respuestaDeserializacionError != nil {
			log.Fatalln("error al decodificar mensaje BSON:", respuestaDeserializacionError)
		}

		artistasRecomendados++
		if strings.HasPrefix(strings.ToLower(respuestaWs.Artista), "n") {
			artistasN++
		}

		for _, c := range respuestaWs.Identificador {
			if unicode.IsLetter(c) {
				letrasIdentificador++
			}
		}

		if respuestaWs.Fin {
			break
		}
	}
	fmt.Printf(`a) %d
b) %d
c) %d
`, artistasRecomendados, artistasN, letrasIdentificador)
}

type MensajeWS struct {
	Identificador string
	Espera        int32
	Fin           bool
}

type RespuestaWs struct {
	Error         bool
	Identificador string
	Artista       string
	Fin           bool
}
