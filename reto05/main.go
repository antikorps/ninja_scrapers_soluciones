package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

/*
a) ¿Cuál es el número total de géneros presentes?
b) ¿Cuál es la suma de los dos años de formación con más artistas?
c) ¿Cuál es el primer año de formación de artistas del que aparecen registros?
*/
func main() {
	var baseUrl string
	flag.StringVar(&baseUrl, "url", "https://ninjascrapers-production.up.railway.app/html/reto05/", "url base de la api para el reto")
	// var espera int
	// flag.IntVar(&espera, "espera", 10, "tiempo de espera entre peticiones")
	var usuario string
	flag.StringVar(&usuario, "usuario", "ninja", "nombre de usuario")
	var password string
	flag.StringVar(&password, "password", "Scraper13%", "password del usuario")

	flag.Parse()

	jar, jarError := cookiejar.New(nil)
	if jarError != nil {
		log.Fatalln("error creando el cookiejar", jarError)
	}

	cliente := http.Client{
		Timeout: 7 * time.Second,
		Jar:     jar,
	}

	// Autentificacion
	var autentificacionSolicitud = strings.NewReader(`usuario=` + usuario + `&password=` + password)
	autententificacionUrl := baseUrl + "login"
	peticionAutentificacion, peticionAutentificacionError := http.NewRequest("POST", autententificacionUrl, autentificacionSolicitud)
	if peticionAutentificacionError != nil {
		log.Fatalln("error al preparar la petición de autentificacion", peticionAutentificacionError)
	}
	incorporarCabeceras(peticionAutentificacion)

	respuestaAutentificacion, respuestaAutentificacionError := cliente.Do(peticionAutentificacion)
	if respuestaAutentificacionError != nil {
		log.Fatalln("error en la respuesta de autentificacion", respuestaAutentificacionError)
	}
	if respuestaAutentificacion.StatusCode != 200 {
		log.Fatalln("error en el status code esperado después de la autentificación:", respuestaAutentificacion.Status)
	}
	defer respuestaAutentificacion.Body.Close()

	// Página de géneros
	html, htmlError := goquery.NewDocumentFromReader(respuestaAutentificacion.Body)
	if htmlError != nil {
		log.Fatalln("error al parsear el html de la página géneros", htmlError)
	}

	var totalGeneros int
	html.Find("summary").Each(func(i int, s *goquery.Selection) {
		totalGeneros++
	})

	// Buscar enlaces
	var artistasEnlaces []string
	html.Find("main a").Each(func(i int, s *goquery.Selection) {
		href, existencia := s.Attr("href")
		if !existencia {
			log.Println("ATENCIÓN: enlace sin href", s.Text())
		}
		artista := strings.TrimPrefix(href, "/html/reto05/")
		artistaUrl := baseUrl + artista
		if contiene(artistasEnlaces, artistaUrl) {
			return
		}
		artistasEnlaces = append(artistasEnlaces, baseUrl+artista)
	})

	var formacionArtistas = make(map[int]int)
	for _, v := range artistasEnlaces {
	realizarPeticion:
		peticionArtista, peticionArtistaError := http.NewRequest("GET", v, nil)
		if peticionArtistaError != nil {
			log.Fatalln("error preparando la petición para el artista", v, peticionArtistaError)
		}
		incorporarCabeceras(peticionArtista)

		respuestaArtista, respuestaArtistaError := cliente.Do(peticionArtista)
		if respuestaArtistaError != nil {
			log.Fatalln("error en la respuesta para el artista", v, respuestaArtistaError)
		}
		defer respuestaArtista.Body.Close()
		if respuestaArtista.StatusCode == 429 {
			aviso, avisoError := io.ReadAll(respuestaArtista.Body)
			if avisoError != nil {
				log.Fatalln("no se ha podido leer la respuesta de un status code 429", avisoError)
			}

			expRegSegundos := regexp.MustCompile(`\d{1,2}`)
			coincidencias := expRegSegundos.FindStringSubmatch(string(aviso))
			if (len(coincidencias)) != 1 {
				log.Fatalln("no se ha podido gestionar la espera tras un status code 429 al no capturar el aviso de los segundos")
			}
			espera, esperaError := strconv.Atoi(coincidencias[0])
			if esperaError != nil {
				log.Fatalln("no se ha podido gestionar la espera tras un status code 429 al no convertir a int los segundos", esperaError)
			}
			time.Sleep(time.Duration(espera+1) * time.Second)
			goto realizarPeticion

		}
		if respuestaArtista.StatusCode != 200 {
			log.Fatalln("error en el status code", respuestaArtista.Status, "para el artista", v)
		}
		defer respuestaArtista.Body.Close()
		htmlArtista, htmlArtistaError := goquery.NewDocumentFromReader(respuestaArtista.Body)
		if htmlArtistaError != nil {
			log.Fatalln("error al parsear el html del artista", v, htmlArtistaError)
		}

		formacion := strings.TrimSpace(htmlArtista.Find("h2").Text())
		fecha, fechaError := strconv.Atoi(formacion)
		if fechaError != nil {
			log.Fatalln("no se ha podido convertir a entero una fecha de formación", fechaError)
		}
		valor, existe := formacionArtistas[fecha]
		if !existe {
			formacionArtistas[fecha] = 1
		} else {
			formacionArtistas[fecha] = valor + 1
		}
	}

	var coleccionFormacion []int

	type artistasFormacion struct {
		formacion int
		numero    int
	}
	var artistasFormacionRecopilacion []artistasFormacion

	for c, v := range formacionArtistas {
		coleccionFormacion = append(coleccionFormacion, c)
		artistasFormacionRecopilacion = append(artistasFormacionRecopilacion, artistasFormacion{
			formacion: c,
			numero:    v,
		})
	}

	sort.Ints(coleccionFormacion)

	sort.Slice(artistasFormacionRecopilacion, func(i, j int) bool {
		return artistasFormacionRecopilacion[i].numero > artistasFormacionRecopilacion[j].numero
	})

	sumaAños := artistasFormacionRecopilacion[0].formacion + artistasFormacionRecopilacion[1].formacion

	fmt.Printf(`a) %d
b) %d
c) %d
`, totalGeneros, sumaAños, coleccionFormacion[0])

}

func contiene(coleccion []string, elemento string) bool {
	for _, v := range coleccion {
		if v == elemento {
			return true
		}
	}
	return false
}

func incorporarCabeceras(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:122.0) Gecko/20100101 Firefox/122.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "es-ES,en;q=0.5")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "no-cache")
}
