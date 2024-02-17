package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

/*
a) ¿Cuántos nombres de artistas empiezan por la letra C?
b) ¿Cuántas grabaciones fueron editadas en el año 1995?
c) ¿Cuántos veces aparece el género 'electronic' entre todos los artistas?
*/
func main() {
	var baseUrl string
	flag.StringVar(&baseUrl, "url", "https://ninjascrapers-production.up.railway.app/html/reto01/", "url base del reto")
	flag.Parse()

	cliente := http.Client{
		Timeout: 7 * time.Second,
	}

	// Obtener sitemap
	mapSitioUrl := baseUrl + "sitemap.xml"
	mapaSitioPeticion, mapaSitioPeticionError := http.NewRequest("GET", mapSitioUrl, nil)
	if mapaSitioPeticionError != nil {
		log.Fatalln("ha fallado la mapaSitioPeticion", mapaSitioPeticionError)
	}
	mapaSitioRespuesta, mapaSitioRespuestaError := cliente.Do(mapaSitioPeticion)
	if mapaSitioRespuestaError != nil {
		log.Fatalln("ha fallado mapaSitioRespuesta", mapaSitioRespuestaError)
	}
	if mapaSitioRespuesta.StatusCode != 200 {
		log.Fatalln("al consultar el mapa del sitio se ha obtenido un status code incorrecto", mapaSitioRespuesta.Status)
	}
	defer mapaSitioRespuesta.Body.Close()
	mapaSitioContenido, mapaSitioContenidoError := io.ReadAll(mapaSitioRespuesta.Body)
	if mapaSitioContenidoError != nil {
		log.Fatalln("no se ha podido leer el contenido del mapa del sitio", mapaSitioContenidoError)
	}

	var mapaSitio MapaSitio
	mapaSitioDeserializacionError := xml.Unmarshal(mapaSitioContenido, &mapaSitio)
	if mapaSitioDeserializacionError != nil {
		log.Fatalln("ha fallado la deserializacion del mapa del sitio")
	}

	// Recolección urls
	var analisisUrls []string
	for _, v := range mapaSitio.URL {
		url := v.Loc
		if strings.HasPrefix(url, "/html/reto01/posts/") || strings.HasSuffix(url, "/posts/") {
			entradaUrl := strings.TrimPrefix(url, "/html/reto01/")
			analisisUrls = append(analisisUrls, baseUrl+entradaUrl)
		}
	}

	// Análisis
	var artistas int
	var publicacion int
	var generos int

	lotes := dividirEnLotes(analisisUrls, 5)
	for _, lote := range lotes {
		canal := make(chan AnalisisEntrada)
		for _, v := range lote {
			wg.Add(1)
			go analizarEntradaBlog(&cliente, v, canal)
		}

		go func() {
			wg.Wait()
			close(canal)
		}()

		for v := range canal {
			if v.error {
				log.Println("ERROR:", v.mensajeError)
			}
			artistas += v.artistaC
			generos += v.generoElectronic
			publicacion += v.publicacion1995
		}
	}

	fmt.Printf(`a) %d
b) %d
c) %d
`, artistas, publicacion, generos)

}

var wg sync.WaitGroup

type MapaSitio struct {
	URL []struct {
		Text string `xml:",chardata"`
		Loc  string `xml:"loc"`
	} `xml:"url"`
}

type AnalisisEntrada struct {
	error            bool
	mensajeError     string
	artistaC         int
	publicacion1995  int
	generoElectronic int
}

func dividirEnLotes[T any](coleccion []T, longitud int) (lotes [][]T) {
	for longitud < len(coleccion) {
		coleccion, lotes = coleccion[longitud:], append(lotes, coleccion[0:longitud:longitud])
	}
	return append(lotes, coleccion)
}

func analizarEntradaBlog(cliente *http.Client, url string, canal chan AnalisisEntrada) {
	defer wg.Done()
	peticionEntrada, peticionEntradaError := http.NewRequest("GET", url, nil)
	if peticionEntradaError != nil {
		canal <- AnalisisEntrada{
			error:        true,
			mensajeError: fmt.Sprintf("fallo preparando petición en %v: %v", url, peticionEntradaError),
		}
		return
	}

	respuestaEntrada, respuestaEntradaError := cliente.Do(peticionEntrada)
	if respuestaEntradaError != nil {
		canal <- AnalisisEntrada{
			error:        true,
			mensajeError: fmt.Sprintf("fallo respuesta petición en %v: %v", url, respuestaEntradaError),
		}
		return
	}

	if respuestaEntrada.StatusCode != 200 {
		canal <- AnalisisEntrada{
			error:        true,
			mensajeError: fmt.Sprintf("status code incorrecto en %v: %v", url, respuestaEntrada.Status),
		}
		return
	}
	defer respuestaEntrada.Body.Close()

	html, htmlError := goquery.NewDocumentFromReader(respuestaEntrada.Body)
	if htmlError != nil {
		canal <- AnalisisEntrada{
			error:        true,
			mensajeError: fmt.Sprintf("error en el parseo html en %v: %v", url, htmlError),
		}
		return
	}

	// Artista
	var artista int
	a := html.Find("h1.artista").Text()
	if strings.HasPrefix(strings.ToLower(a), "c") {
		artista++
	}

	var publicaciones int
	html.Find("table tr td:nth-of-type(2)").Each(func(i int, s *goquery.Selection) {
		año := s.Text()
		if año == "1995" {
			publicaciones++
		}
	})

	var generos int
	html.Find(".genero").Each(func(i int, s *goquery.Selection) {
		genero := s.Text()
		if strings.ToLower(genero) == "electronic" {
			generos++
		}
	})

	canal <- AnalisisEntrada{
		artistaC:         artista,
		generoElectronic: generos,
		publicacion1995:  publicaciones,
	}

}
