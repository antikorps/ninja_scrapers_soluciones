package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
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

	expRegData := regexp.MustCompile(`.*?function Gi\(i,o,t\)\{let a=(.*?),c=0,n=a`)

	coincidencias = expRegData.FindStringSubmatch(string(scriptContenido))

	if len(coincidencias) != 2 {
		log.Fatalln("no se ha podido capturar el grupo con la información en JSON")
	}

	data := coincidencias[1]

	// Años de formacion
	var artistasFundados19801985 int
	expRegBusqueda := regexp.MustCompile(`formacion:(\d{4}),`)
	coincidenciasBusquedas := expRegBusqueda.FindAllStringSubmatch(data, -1)

	for _, v := range coincidenciasBusquedas {
		if len(v) != 2 {
			fmt.Println("una coincidencia de formación no tiene grupo capturado, probablemente el script de resultados incorrectos", v)
			continue
		}
		año, añoError := strconv.Atoi(v[1])
		if añoError != nil {
			fmt.Println("un grupo capturado no es un número, probablemente el script de resultados incorrectos", v)
			continue
		}
		if año >= 1980 && año <= 1985 {
			artistasFundados19801985++
		}
	}

	// Longitud del título más corto
	tituloMasCorto := 100000000
	expRegBusqueda = regexp.MustCompile(`"titulo":.*?"(.*?)"}`)
	coincidenciasBusquedas = expRegBusqueda.FindAllStringSubmatch(data, -1)

	for _, v := range coincidenciasBusquedas {
		if len(v) != 2 {
			fmt.Println("una coincidencia de tituloMasCorto no tiene grupo capturado, probablemente el script de resultados incorrectos", v)
			continue
		}
		// Eliminar caracteres escapados
		titulo := strings.Replace(v[1], `\\`, "", -1)
		longitudTitulo := len([]rune(titulo))
		if longitudTitulo < tituloMasCorto {
			tituloMasCorto = longitudTitulo
		}
	}

	// Suma números identificadores
	mayorSumaIdentificadores := 0
	expRegBusqueda = regexp.MustCompile(`identificador:"(.*?)"`)
	coincidenciasBusquedas = expRegBusqueda.FindAllStringSubmatch(data, -1)
	for _, v := range coincidenciasBusquedas {
		if len(v) != 2 {
			fmt.Println("una coincidencia de tituloMasCorto no tiene grupo capturado, probablemente el script de resultados incorrectos", v)
			continue
		}
		var sumaIdentificador int
		for _, c := range v[1] {
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
