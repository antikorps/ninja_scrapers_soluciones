import requests
import time

URL_BASE = "https://ninjascrapers-production.up.railway.app/api/reto03"
CABECERAS = {
    "User-Agent":"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:122.0) Gecko/20100101 Firefox/122.0",
	"Accept-Language": "es-ES:en;q=0.5",
	"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
	"DNT": "1",
	"Connection": "keep-alive",
	"Sec-Fetch-Dest": "empty",
	"Sec-Fetch-Mode": "cors",
	"Sec-Fetch-Site": "same-origin",
	"Pragma": "no-cache",
	"Cache-Control": "no-cache",
	"Content-Length": "0"
}

def main():
    """a) ¿Cual es el número total de registros que la página muestra? 
b) ¿Cuál es el número total de grabaciones que aparecen en los registros?
c) ¿Cuál es el número total de géneros si se suman todos los que aparezcan (sim importar que estén repetidos)?
    """
    sesion = requests.Session()

    offset = 0
    total_registros = 0
    total_grabaciones = 0
    total_generos = 0

    while True:
        parametros = {
			"t":      int(time.time()),
			"id":     "discografía completa",
			"num":    "25",
			"offset": offset
		}
        r = sesion.post(URL_BASE, headers=CABECERAS, params=parametros)
        if r.status_code != 200:
            mensaje_error = f"status code no esperado {r.status_code} procesando el offset {offset}"
            raise Exception(mensaje_error)
        respuesta = r.json()

        for registro in respuesta:
            total_registros += 1
            total_grabaciones += len(registro["discografia"])
            total_generos += len(registro["generos"])
        
        if len(respuesta) < 25:
            break

        offset += 25

    print(f"""a) {total_registros}
b) {total_grabaciones}
c) {total_generos}""")


main()
