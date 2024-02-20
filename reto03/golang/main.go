package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

/*
a) ¿Cual es el número total de registros que la página muestra?
b) ¿Cuál es el número total de grabaciones que aparecen en los registros?
c) ¿Cuál es el número total de géneros si se suman todos los que aparezcan (sim importar que estén repetidos)?
*/

func main() {
	var baseUrl string
	flag.StringVar(&baseUrl, "url", "https://ninjascrapers-production.up.railway.app/api/reto03", "url base de la api para el reto")
	flag.Parse()

	cliente := http.Client{
		Timeout: 7 * time.Second,
	}

	var offset int
	var totalRegistros int
	var totalGrabaciones int
	var totalGeneros int
	for {
		peticionApi, peticionApiError := http.NewRequest("POST", baseUrl, nil)
		if peticionApiError != nil {
			log.Fatalln("error en la petición a la api", peticionApiError, offset)
		}
		incorporarCabeceras(peticionApi)
		t := time.Now().Unix()
		parametros := map[string]string{
			"t":      fmt.Sprint(t),
			"id":     "discografía completa",
			"num":    "25",
			"offset": fmt.Sprint(offset),
		}
		urlQuery := peticionApi.URL.Query()
		for c, v := range parametros {
			urlQuery.Add(c, v)
		}
		peticionApi.URL.RawQuery = urlQuery.Encode()

		respuestaApi, respuestaApiError := cliente.Do(peticionApi)
		if respuestaApiError != nil {
			log.Fatalln("error en la respuesta de la api", respuestaApiError, offset)
		}
		if respuestaApi.StatusCode != 200 {
			log.Fatalln("error en el status code de la api", respuestaApi.Status, offset)
		}
		var r R

		errorDeserializacion := json.NewDecoder(respuestaApi.Body).Decode(&r)
		if errorDeserializacion != nil {
			log.Fatalln("error en la deserialización de la respuesta de la api", errorDeserializacion)
		}

		for _, v := range r {
			totalRegistros++
			totalGrabaciones += len(v.Discografia)
			totalGeneros += len(v.Generos)
		}
		if len(r) < 25 {
			break
		}
		offset += 25
	}

	fmt.Printf(`a) %d
b) %d
c) %d
`, totalRegistros, totalGrabaciones, totalGeneros)
}

type R []struct {
	Artista     string `json:"artista"`
	Discografia []struct {
		Publicacion int    `json:"publicacion"`
		Titulo      string `json:"titulo"`
	} `json:"discografia"`
	Generos []string `json:"generos"`
}

func incorporarCabeceras(peticion *http.Request) {
	peticion.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:122.0) Gecko/20100101 Firefox/122.0")
	peticion.Header.Set("Accept", "*/*")
	peticion.Header.Set("Accept-Language", "es-ES,en;q=0.5")
	peticion.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	peticion.Header.Set("DNT", "1")
	peticion.Header.Set("Connection", "keep-alive")
	peticion.Header.Set("Sec-Fetch-Dest", "empty")
	peticion.Header.Set("Sec-Fetch-Mode", "cors")
	peticion.Header.Set("Sec-Fetch-Site", "same-origin")
	peticion.Header.Set("Pragma", "no-cache")
	peticion.Header.Set("Cache-Control", "no-cache")
	peticion.Header.Set("Content-Length", "0")
}
