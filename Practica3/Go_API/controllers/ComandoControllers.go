package Controllers

import (
	Models "Backend/models"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"os/user"
	"io"
	"strconv"
	"bufio"
	"os"
)
/**
 Este método ejecuta los comandos utilizando os/exec
 @param {string} comando  Es el comando a ejecutar
 @return {string}  Salida del comando
*/
func CMD(comando string) (bytes.Buffer, string, error) {
	// Crea dos búferes de bytes para capturar la salida y los errores
	var salida bytes.Buffer
	var errors bytes.Buffer
	
	// Ejecuta el comando proporcionado en una shell de Bash
	cmd := exec.Command("bash", "-c", comando)
	cmd.Stdout = &salida
	cmd.Stderr = &errors
	err := cmd.Run()
	
	// Devuelve la salida capturada, los errores como cadena de texto y cualquier error encontrado
	return salida, errors.String(), err
}

/**
    Retorna el json creado por el módulo de kernel en el directorio /proc de los procesos del sistema,
	donde además se convierte el id del usuario propietario del proceso al nombre real del usuario

	@param: null

	@return: res (json): json con la información obtenida del archivo escrito por el módulo de kernel
*/
func RequestPrincipal() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// se definen los headers necesarios para realizar las peticiones
		rw.Header().Set("Content-Type", "application/json")
		
		// Leer el archivo JSON en el directorio /proc
		salida, _, verificar := CMD("cat /proc/cpu_grupo18")
		if verificar != nil {
			log.Printf("error: %v\n", verificar)
			return
		}
		
		// Parsear el JSON a una estructura CPUDATAJSON
		var dataJson Models.CPUDATAJSON
		if err := json.Unmarshal(salida.Bytes(), &dataJson); err != nil {
			// se imprime error y se retorna 
			log.Printf("error: %v\n", err)
			return
		}
		
		
		// Recorrer los datos y convertir los IDs de usuario a nombres de usuario
		for i := range dataJson.DATA {
			// obtencion del usuario correspondiente a un ID de usuario dado
			username, err := getUsername(dataJson.DATA[i].USUARIO)
				if err != nil {
					log.Printf("Error al obtener el nombre de usuario: %v\n", err)
					http.Error(rw, "Error interno del servidor", http.StatusInternalServerError)
					return
				}
			// se guarda el nuevo nombre de usuario dentro del objeto
			dataJson.DATA[i].USUARIO = username
			// se recorren los procesos hijos para buscar el nombre de usuario al que pertenece
			// cada uno de ellos, tambien se actualiza el valor del objeto con el nuevo nombre de
			// usuario
			for j := range dataJson.DATA[i].PROCESOSHIJOS {
				username, err := getUsername(dataJson.DATA[i].PROCESOSHIJOS[j].USUARIO)
				if err != nil {
					log.Printf("Error al obtener el nombre de usuario: %v\n", err)
					http.Error(rw, "Error interno del servidor", http.StatusInternalServerError)
					return
				}
				// actualizacion del nombre de usuario
				dataJson.DATA[i].PROCESOSHIJOS[j].USUARIO = username

			}
		}
		
		// Codificar la estructura actualizada y enviarla en la respuesta
		if err := json.NewEncoder(rw).Encode(dataJson); err != nil {
			log.Printf("error: %v\n", err)
			return
		}
	}
}


/**
 * Este método obtiene el nombre del usuario.
 * @param  {string} userID Es el id obtenido de los módulos.
 * @return {string} El nombre del usuario
 */
func getUsername(userID string) (string, error) {
	//Verificación del nombre de usuario perteneciente a un ID X
	u, err := user.LookupId(userID)  
	if err != nil {
		return "", err
	}
	return u.Username, nil
}


/*
    Se ejecuta al ingresar al endpoint /Kill

    @params {int} pid id del proceso que se quiere eliminar
	@returns {json} confirmación de eliminación del proceso
*/
func RequestKill() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Verificamos que el endpoint sea el correcto
		if r.URL.Path != "/Kill" {
			http.NotFound(rw, r)
			return
		}
		// se valida el metodo necesario para servir la peticion.
		if  r.Method =="GET" {
			//Obtención del parámetro pid en ruta
			id := r.URL.Query().Get("pid")
			id = strings.TrimSuffix(id, "/")
			//Ejecución del comando para la eliminación del proceso con el pid recibido
			  _, _, verificar := CMD("sudo kill -9 " + id)

			//Verificación de eliminación
			if verificar != nil {
				log.Printf("error: %v\n", verificar)
			} else {
				fmt.Println("Eliminando Proceso: " + id)
			} 
		}else{
				//sino se reconoce el metodo se alerta el estado del error.
				rw.WriteHeader(http.StatusNotImplemented)
				rw.Write([]byte(http.StatusText(http.StatusNotImplemented)))
		}
		
	}
}

/*
    Retorna el json creado por el módulo de kernel en el directorio /proc de la memoria ram utilizada
	por el sistema
	
	@returns (json) objeto con la memoria ram utilizada por el sistema escrito por el módulo 
			de kernel
*/
func RequestMemory() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		//Lectura del archivo en el directorio /proc
		salida, _, verificar := CMD("cat /proc/mem_grupo18")

		//Verificando la correcta lectura del archivo
		if verificar != nil {
			log.Printf("error: %v\n", verificar)
		} else {
			var dataJson Models.DATAJSONMEMORY // se declara el objeto json que contiene datos de la memoria
			json.Unmarshal(salida.Bytes(), &dataJson) //se hace UnMarshall de json a objeto
      		rw.Header().Set("Content-Type", "application/json") //Agregando headers para que sea json
			json.NewEncoder(rw).Encode(dataJson) // Enviamos la información por el responseWriter
		}
	}
}

/*
*   Retorna un mensaje para probar el funcionamiento del servidor
*
* @Param: null
* @return (cadena): cadena que informa que funciona la ruta /
*/
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API GO!\n"))
}

/*
*  Retorna las ubicaciones donde se está ejecutando el proceso que se solicita
*
* @Param (int): Identificador del proceso del que se requiere obtener su administración de memoria
* @return {json}: lista de jsons que poseen la información sobre la administración de memoria del proceso requerido
*/
func RequestMaps() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Verificamos que la ruta sea correcta
		if r.URL.Path != "/maps" {
			http.NotFound(rw, r) // si no encuentra la ruta retorna un No encontrado
			return
		}
		if  r.Method =="GET" {
			//Obtención del parámetro pid en ruta
			id := r.URL.Query().Get("pid")
			id = strings.TrimSuffix(id, "/") // Limpiamos el parámetro
			// casteando a entero
			num, err := strconv.Atoi(id)
			if err != nil {
				fmt.Println("Error al convertir el string a int:", err)
				return
			}
			var memorys []Models.Memory  // variable representa al Mapa de datos de la Memoria
			//Ejecución del comando para la eliminación del proceso con el pid recibido
			salida, verificar := ObtenerDatosMaps(num)
			if verificar != nil {
				log.Printf("error: %v\n", verificar)
				return 
			} 
			//! ------ Obtener RSS y Size del archivo /proc/pid/smaps  -------
			salidasmaps, smapsErr := ObtenerRSS(num)
			if smapsErr != nil {
				log.Printf("Error en la lectura de SMAPS: %v", err)
				return
			}
			memory := Models.Memory{
				Arr1:   salida,
				Arr2:   salidasmaps,
			}
			//Añadimos el objeto a la lista
			memorys = append(memorys, memory)
			// Agregando cabecera para que se formatee a json
      		rw.Header().Set("Content-Type", "application/json")
			json.NewEncoder(rw).Encode(memorys) //Enviamos salida por el responseWriter
		}else{
			//sino se reconoce el metodo se alerta el estado del error.
			rw.WriteHeader(http.StatusNotImplemented)
			rw.Write([]byte(http.StatusText(http.StatusNotImplemented)))
		}
	}
}

/*
*  Retorna las ubicaciones donde se está ejecutando el proceso que se solicita
*
* @Param: (int): Identificador del proceso del que se requiere obtener su administración de memoria
* @return: {MemoryMap}: lista de jsons que poseen la información sobre la administración de memoria del proceso requerido
*/
func ObtenerDatosMaps(pid int) ([]Models.MemoryMap, error) {
	// Leemos el contenido de /proc/pid/maps para obtener la administración de su memoria
	mapsPath := fmt.Sprintf("/proc/%d/maps", pid)
	file, err := os.Open(mapsPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var memoryMaps []Models.MemoryMap  // variable representa al Mapa de datos de la Memoria

	reader := bufio.NewReader(file)
	// se recorre ciclicamente la lectura y la ejecucion para la verificacion de procesos y memoria RAM
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}

		if line == "" || err == io.EOF {
			break
		}

		fields := strings.Fields(line) //dividir la línea de texto en campos individuales.
		//Se verifica si el slice fields tiene al menos 6 elementos antes de continuar con la extracción
		if len(fields) >= 6 {
			addressRange := fields[0] //rango de direcciones.
			permissions := fields[1] //los permisos asociados al rango de direcciones.
			device := fields[3] //el dispositivo asociado al rango de direcciones.
			filePath := fields[5] //ruta de archivo asociada al rango de direcciones.

			// Obtener el tamaño del rango de direcciones
			rangeFields := strings.Split(addressRange, "-")
			if len(rangeFields) != 2 {
				continue
			}
			// Casteando la dirección de inicio a enteros
			start, err := strconv.ParseUint(rangeFields[0], 16, 64)
			if err != nil {
				log.Printf("Error al parsear la dirección de inicio: %v", err)
				continue
			}
			// Casteando la dirección de fin a enteros
			end, err := strconv.ParseUint(rangeFields[1], 16, 64)
			if err != nil {
				log.Printf("Error al parsear la dirección de fin: %v", err)
				continue
			}
			// Cálculo del tamaño de segmento
			size := end - start
			
			memoryMap := Models.MemoryMap{
				Direccion:   addressRange,
				Tamanio:     size,
				Permisos:    permissions,
				Dispositivo: device,
				Archivo:     filePath,
			}
			//Añadimos el objeto a la lista
			memoryMaps = append(memoryMaps, memoryMap)
		}
	}
	return memoryMaps, nil
}

// sudo cat smaps | grep 'Rss:'
// sudo cat smaps | grep 'Size:'
func ObtenerRSS(pid int) ([]Models.MemResVirtual, error) {
	var memResVirtuals []Models.MemResVirtual  // variable representa al Mapa de datos de la Memoria
	salida, _, verificar := CMD("cat /proc/"+strconv.Itoa(pid)+"/smaps | grep 'Rss:'")
	if verificar != nil {
		log.Printf("error: %v\n", verificar)
		return memResVirtuals, nil
	}

	salida2, _, verificar2 := CMD("cat /proc/"+strconv.Itoa(pid)+"/smaps | grep '^Size:'")
	if verificar2 != nil {
		log.Printf("error: %v\n", verificar2)
		return memResVirtuals, nil
	}

	lines := strings.Split(salida.String(), "\n")
	totalRss := 0
	lines2 := strings.Split(salida2.String(), "\n")
	totalSize := 0
	indx := 0

	for _, line := range lines {
		// Extraer el valor numérico de la línea
		valueStr := strings.TrimSpace(strings.TrimPrefix(line, "Rss:"))
		valueStr2 := strings.TrimSpace(strings.TrimPrefix(lines2[indx], "Size:"))
		va := strings.Split(valueStr, " ")
		va2 := strings.Split(valueStr2, " ")
		value, err := strconv.Atoi(va[0])
		value2, err := strconv.Atoi(va2[0])
		
		if err != nil {
			continue
		}

		// Sumar el valor al total
		totalRss += value
		totalSize += value2
		indx += 1
		
		memResVirtual := Models.MemResVirtual{			
			Residente:   value,
			Virtual:     value2,
		}
		//Añadimos el objeto a la lista
		memResVirtuals = append(memResVirtuals, memResVirtual)
	}
	fmt.Printf("Rss total: %d\n", totalRss)
	fmt.Printf("Size total: %d\n", totalSize)
	return memResVirtuals, nil
	// smapsPath := fmt.Sprintf("/proc/%d/smaps", pid)
	// file, err := os.Open(smapsPath)

	// var memResVirtuals []Models.MemResVirtual  // variable representa al Mapa de datos de la Memoria
	// if err != nil {
	// 	return memResVirtuals, err
	// }
	// defer file.Close()

	// reader := bufio.NewReader(file)
	// var rss uint64
	// var sizevm uint64
	// contador := 0
	// contador2 := 0

	// for {
		
	// 	line, err := reader.ReadString('\n')
	// 	if err != nil && err != io.EOF {
	// 		return memResVirtuals, err
	// 	}

	// 	if line == "" || err == io.EOF {
	// 		break
	// 	}

	// 	if strings.HasPrefix(line, "Rss:") {
	// 		contador2=contador2+1
			
	// 		fields := strings.Fields(line)
	// 		if len(fields) >= 2 {
	// 			sizeKb, err := strconv.ParseUint(fields[1], 10, 64)
	// 			if err != nil {
	// 				return memResVirtuals, err
	// 			}
	// 			rss += sizeKb / 1024
	// 		}
			
	// 	}
	// 	if strings.HasPrefix(line, "Size:") {
	// 		contador=contador+1
			
	// 		fields := strings.Fields(line)
	// 		if len(fields) >= 2 {
	// 			sizeKb, err := strconv.ParseUint(fields[1], 10, 64)
	// 			if err != nil {
	// 				return memResVirtuals, err
	// 			}
	// 			sizevm += sizeKb / 1024
	// 		}
	// 		memResVirtual := Models.MemResVirtual{			
	// 			Residente:   rss,
	// 			Virtual:     sizevm,
	// 		}
	// 		//Añadimos el objeto a la lista
	// 		memResVirtuals = append(memResVirtuals, memResVirtual)
	// 	}	
		
	// }
	// fmt.Println(contador)

	// return memResVirtuals, nil
}


// att. el Grupo 18, el mejor!! el grupo más sistémico 😂