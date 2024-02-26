/*
a) ¿Cuántos artistas empiezan por la letra R?
b) ¿Cuántas grabaciones se publicaron en 1995?
c) El identificador de cada artista es un código alfanumérico. Si cuentas cada caracter como una cifra, ¿cuál es la suma de todos los caracteres numéricos encontrontrados?
*/

const sleep = function (ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
};

async function main() {
    const identificadores = []
    const $boton = document.querySelector("button")

    let artistasR = 0
    let grabaciones1995 = 0
    let identificadorNumeros = 0

    while (identificadores.length < 167) {
        const identificador = document.querySelector(".contenedor-identificador").textContent.trim()
        if (!identificadores.includes(identificador)) {
            identificadores.push(identificador)

            const artista = document.querySelector("h2").textContent.trim().toLowerCase()
            if (artista.startsWith("r")) {
                artistasR++
            }

            const $publicaciones = document.querySelectorAll(".publicacion")
            for (const $publicacion of $publicaciones) {
                const fecha = $publicacion.textContent.trim().replace("(", "").replace(")", "")
                if (fecha == "1995") {
                    grabaciones1995++
                }
            }
        }
        await sleep(1500)
        $boton.click()
    }

    for (const identificador of identificadores) {
        for (const caracter of identificador) {
            if (isNaN(caracter)) {
                continue
            }
            const cifra = parseInt(caracter)
            identificadorNumeros += cifra
        }
    }

    console.log(`a) ${artistasR}
b) ${grabaciones1995}
c) ${identificadorNumeros}    
`)

}

main()
