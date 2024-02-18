package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dop251/goja"
)

/*
a) ¿Cuántos artistas se fundaron entre 1980 y 1985 (ambos inclusive)?
b) ¿Cuál es la longitud de caracteres del título más corto de una grabación?
c) El identificador de cada artista es un código alfanumérico. Si se suman todos los caracteres numéricos del mismo, ¿cuál es el mayor resultado?
*/
func main() {
	var baseUrl string
	flag.StringVar(&baseUrl, "url", "https://ninjascrapers-production.up.railway.app/html/reto02/", "url base del reto")
	flag.Parse()

	cliente := http.Client{
		Timeout: 7 * time.Second,
	}

	indexPeticion, indexPeticionError := http.NewRequest("GET", baseUrl, nil)
	if indexPeticionError != nil {
		log.Fatalln("error en la petición para el index", indexPeticionError)
	}
	indexRespuesta, indexRespuestaError := cliente.Do(indexPeticion)
	if indexRespuestaError != nil {
		log.Fatalln("error en la respuesta de index", indexRespuestaError)
	}
	defer indexRespuesta.Body.Close()
	if indexRespuesta.StatusCode != 200 {
		log.Fatalln("error en la respuesta a index al recibir un status code incorrecto", indexRespuesta.Status)
	}

	indexContenido, indexContenidoError := io.ReadAll(indexRespuesta.Body)
	if indexContenidoError != nil {
		log.Fatalln("error al no poder recuperar el contenido la respuesta de index", indexContenidoError)
	}

	expRegScript := regexp.MustCompile(`.*<script.*?src="(.*?)".*`)

	coincidencias := expRegScript.FindStringSubmatch(string(indexContenido))
	if len(coincidencias) != 2 {
		log.Fatalln("error al no poder capturar el src del script")
	}
	scriptUrl := strings.Replace(baseUrl, "/html/reto02/", coincidencias[1], 1)

	scriptPeticion, scriptPeticionError := http.NewRequest("GET", scriptUrl, nil)
	if scriptPeticionError != nil {
		log.Fatalln("error en la petición para el script", scriptPeticionError)
	}

	scriptRespuesta, scriptRespuestaError := cliente.Do(scriptPeticion)
	if scriptRespuestaError != nil {
		log.Fatalln("error en la respuesta del script", scriptRespuestaError)
	}
	defer scriptRespuesta.Body.Close()
	if scriptRespuesta.StatusCode != 200 {
		log.Fatalln("error en la respuesta del script al recibir un status code incorrecto", scriptRespuesta.Status)
	}

	scriptContenido, scriptContenidoError := io.ReadAll(scriptRespuesta.Body)
	if scriptContenidoError != nil {
		log.Fatalln("error al intentar recuperar el contenido del script", scriptContenidoError)
	}

	expRegVariable := regexp.MustCompile(`Gi\(i,o,t\){(let a=.*?),c=0`)
	coincidencias = expRegVariable.FindStringSubmatch((string(scriptContenido)))
	if len(coincidencias) != 2 {
		log.Fatalln("no se ha podido capturar el valor de la variable a en el script")
	}

	gojaRuntime := goja.New()

	codigoJS := fmt.Sprintf("function serializar(){%v;return JSON.stringify(a)}", coincidencias[1])

	_, ejecucionJSError := gojaRuntime.RunString(codigoJS)
	if ejecucionJSError != nil {
		log.Fatalln("error en la ejecución del código js capturado", ejecucionJSError)
	}
	var serializar func() string
	exportarJSError := gojaRuntime.ExportTo(gojaRuntime.Get("serializar"), &serializar)
	if exportarJSError != nil {
		log.Fatalln("error en la exportación de la función serializar", exportarJSError)
	}

	datosSerializados := serializar()

	var data Data
	errorDeserializacion := json.Unmarshal([]byte(datosSerializados), &data)
	if errorDeserializacion != nil {
		log.Fatalln("no se ha podido deserializar en data la ejecución de la función en JS", errorDeserializacion)
	}

	var artistasFundados19801985 int
	tituloMasCorto := 1_000_000
	mayorSumaIdentificadores := 0
	for _, registro := range data {
		if registro.Formacion >= 1980 && registro.Formacion <= 1985 {
			artistasFundados19801985++
		}

		var discografia Discografia
		errorDeserializacionDiscografia := json.Unmarshal([]byte(registro.Discografia), &discografia)
		if errorDeserializacionDiscografia != nil {
			log.Fatalln("no se ha podido deserializar la discografía de", registro.Artista, errorDeserializacionDiscografia)
		}

		for _, disco := range discografia {
			longitudTitulo := len(disco.Titulo)
			if longitudTitulo < tituloMasCorto {
				tituloMasCorto = longitudTitulo
			}
		}

		var sumaIdentificador int
		for _, c := range registro.Identificador {
			cifra, cifraError := strconv.Atoi(string(c))
			if cifraError == nil {
				sumaIdentificador += cifra
			}
		}
		if sumaIdentificador > mayorSumaIdentificadores {
			mayorSumaIdentificadores = sumaIdentificador
		}
	}

	fmt.Printf(`a) %d
b) %d
c) %d
`, artistasFundados19801985, tituloMasCorto, mayorSumaIdentificadores)

}

type Data []struct {
	Identificador string `json:"identificador"`
	Artista       string `json:"artista"`
	Formacion     int    `json:"formacion"`
	Discografia   string `json:"discografia"`
	Generos       string `json:"generos"`
}

type Discografia []struct {
	Publicacion int    `json:"publicacion"`
	Titulo      string `json:"titulo"`
}
