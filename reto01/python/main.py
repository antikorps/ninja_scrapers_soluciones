import requests
from bs4 import BeautifulSoup
import concurrent.futures
from dataclasses import dataclass

URL_BASE = "https://ninjascrapers-production.up.railway.app/html/reto01/"
URL_SITEMAP = URL_BASE + "sitemap.xml"

@dataclass
class Entrada:
    artista: int
    grabaciones: int
    generos: int

def analizar_entrada_blog(sesion: requests.Session, url: str):
    r = sesion.get(url)
    if r.status_code != 200:
        print(f"status code incorrecto {r.status_code} consultando {url}")
        return None
   
    soup = BeautifulSoup(r.text, "html.parser")
    artista = soup.select_one("h1.artista")
    if artista == None:
        print(f"no se ha podido recuperar el artista en {url}")
        return None
    art = 0
    if artista.text.strip().lower().startswith("c"):
        art += 1
    
    publicaciones = soup.select("table tr td:nth-of-type(2)")
    grabaciones = 0
    for p in publicaciones:
        if p.text.strip() == "1995":
            grabaciones += 1

    generos = soup.select(".genero")
    electronic = 0
    for g in generos:
        if g.text.strip() == "electronic":
            electronic += 1

    return Entrada(art, grabaciones, electronic)

# a) ¿Cuántos nombres de artistas empiezan por la letra C?
# b) ¿Cuántas grabaciones fueron editadas en el año 1995?
# c) ¿Cuántos veces aparece el género 'electronic' entre todos los artistas?

def main():
    sesion = requests.Session()

    r = sesion.get(URL_SITEMAP)
    if r.status_code != 200:
        raise Exception(f"status code incorrecto consultando el sitemap {r.status_code}")
    soup = BeautifulSoup(r.text, "xml")
    locs = soup.find_all("loc")

    urls = []
    for loc in locs:
        url = loc.text
        if url.startswith("/html/reto01/posts/") == False or url.endswith("posts/"):
            continue
        urls.append(URL_BASE + url.replace("/html/reto01/", ""))

    artistas = 0
    grabaciones = 0
    electronic = 0

    with concurrent.futures.ThreadPoolExecutor(max_workers=5) as manejador:
        futuros = []
        for url in urls:
            futuros.append(manejador.submit(analizar_entrada_blog, sesion=sesion, url=url))
        for futuro in concurrent.futures.as_completed(futuros):
            resultado: Entrada = futuro.result()
            if resultado == None:
                continue
            artistas += resultado.artista
            grabaciones += resultado.grabaciones
            electronic += resultado.generos

    print(f"""a) {artistas}
b) {grabaciones}
c) {electronic}""")
    

main()