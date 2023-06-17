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
)

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

'''
    Retorna el json creado por el módulo de kernel en el directorio /proc de los procesos del sistema,
	donde además se convierte el id del usuario propietario del proceso al nombre real del usuario

        Parameters:
			-

        Returns:
            res (json): json con la información obtenida del archivo escrito por el módulo de kernel
'''
func RequestPrincipal() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
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
			log.Printf("error: %v\n", err)
			return
		}
		
		
		// Recorrer los datos y convertir los IDs de usuario a nombres de usuario
		for i := range dataJson.DATA {

			username, err := getUsername(dataJson.DATA[i].USUARIO)
				if err != nil {
					log.Printf("Error al obtener el nombre de usuario: %v\n", err)
					http.Error(rw, "Error interno del servidor", http.StatusInternalServerError)
					return
				}
			dataJson.DATA[i].USUARIO = username
			
			for j := range dataJson.DATA[i].PROCESOSHIJOS {
				username, err := getUsername(dataJson.DATA[i].PROCESOSHIJOS[j].USUARIO)
				if err != nil {
					log.Printf("Error al obtener el nombre de usuario: %v\n", err)
					http.Error(rw, "Error interno del servidor", http.StatusInternalServerError)
					return
				}
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


'''
    Retorna el nombre del usuario en base al id obtenido del propietario del proceso

        Parameters:
			userID (string): Cadena que posee el id del propietario del proceso

        Returns:
            usuario (string): nombre del propietario del proceso
			error (error): variable para verificar si ocurrió un error
'''
func getUsername(userID string) (string, error) {
	u, err := user.LookupId(userID) //Verificación de 
	if err != nil {
		return "", err
	}
	return u.Username, nil
}


'''
    Por medio de la petición 

        Parameters:
			/:pid (int): id del proceso que se quiere eliminar

        Returns:
            elim (string): confirmación de eliminación del proceso
'''
func RequestKill() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/Kill" {
			http.NotFound(rw, r)
			return
		}
		if  r.Method =="GET" {
			//Obtención del parámetro pid en ruta
			id := r.URL.Query().Get("pid")
			id = strings.TrimSuffix(id, "/")
      		fmt.Println(id)
			//Ejecución del comando para la eliminación del proceso con el pid recibido
			  _, _, verificar := CMD("sudo kill -9 " + id)

			//Verificación de eliminación
			if verificar != nil {
				log.Printf("error: %v\n", verificar)
			} else {
				fmt.Println("Eliminando Proceso: " + id)
			} 
    }else{
			rw.WriteHeader(http.StatusNotImplemented)
			rw.Write([]byte(http.StatusText(http.StatusNotImplemented)))
    }
		
	}
}

'''
    Retorna el json creado por el módulo de kernel en el directorio /proc de la memoria ram utilizada
	por el sistema

        Parameters:
			-

        Returns:
            res (json): json con la memoria ram utilizada por el sistema escrito por el módulo 
			de kernel
'''
func RequestMemory() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		//Lectura del archivo en el directorio /proc
		salida, _, verificar := CMD("cat /proc/mem_grupo18")

		//Verificando la correcta lectura del archivo
		if verificar != nil {
			log.Printf("error: %v\n", verificar)
		} else {
			var dataJson Models.DATAJSONMEMORY
			json.Unmarshal(salida.Bytes(), &dataJson) //json a objeto
      		rw.Header().Set("Content-Type", "application/json")
			json.NewEncoder(rw).Encode(dataJson)
		}
	}
}

'''
    Retorna un mensaje para probar el funcionamiento del servidor

        Parameters:
			-

        Returns:
            res (cadena): cadena que informa que funciona la ruta /
'''
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API GO!\n"))
}
