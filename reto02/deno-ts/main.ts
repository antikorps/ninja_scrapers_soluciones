function calcularSumaCifras(identificador: string): number {
    let suma = 0
    for (const caracter of identificador) {
        const cifra = parseInt(caracter)
        if (isNaN(cifra)) {
            continue
        }
        suma += cifra
    }
    return suma
}

async function main() {
    let baseUrl = "https://ninjascrapers-production.up.railway.app/html/reto02/"
    let r = await fetch(baseUrl)
    if (r.status != 200) {
        throw Error(`status code incorrecto ${r.statusText} al consultar la página principal`)
    }

    let res = await r.text()
    const expRegScriptSrc = /script.*?src="(.*?)"/m
    const scriptSrc = res.match(expRegScriptSrc)
    if (scriptSrc == null) {
        throw Error(`la expresión regular para capturar el src no ha tenido coincidencias`)
    }
    if (scriptSrc.length != 2) {
        throw Error(`la expresión regular no tiene una longitud de 2, error capturando grupos`)
    }

    const scriptUrl = baseUrl.replace("/html/reto02/", scriptSrc[1])
    
    r = await fetch(scriptUrl)
    if (r.status != 200) {
        throw Error(`status code incorrecto ${r.statusText} al consultar el script`)
    }
    res = await r.text()

    const expRegDatos = /Gi\(i,o,t\){let a=(.*?),c=0/
    const busquedaDatos = res.match(expRegDatos)
    if (busquedaDatos == null) {
        throw Error ("no")
    }
    const devolverDatos = new Function(`return ${busquedaDatos[1]}`)

    const datos = devolverDatos()
    
    let artistasFundados19801985 = 0
    let tituloMasCorto = 1_000_000
    let mayorSumaIdentificadores = 0
    for (const registro of datos) {
        if (registro.formacion >= 1980 && registro.formacion <= 1985) {
            artistasFundados19801985++
        }
        const grabaciones = JSON.parse(registro.discografia)
        for (const grabacion of grabaciones) {
            const longitudTitulo = grabacion.titulo.length
            if (longitudTitulo < tituloMasCorto) {
                tituloMasCorto = longitudTitulo
            }
        }
        const sumaIdentificadores = calcularSumaCifras(registro.identificador)
        if (sumaIdentificadores > mayorSumaIdentificadores) {
            mayorSumaIdentificadores = sumaIdentificadores
        } 
    }

    console.log(`a) ${artistasFundados19801985}
b) ${tituloMasCorto}
c) ${mayorSumaIdentificadores}`)
}

/*
a) ¿Cuántos artistas se fundaron entre 1980 y 1985 (ambos inclusive)?
b) ¿Cuál es la longitud de caracteres del título más corto de una grabación?
c) El identificador de cada artista es un código alfanumérico. Si se suman todos los caracteres numéricos del mismo, ¿cuál es el mayor resultado?
*/
main()