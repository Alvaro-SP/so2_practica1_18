package Controllers

import (
	Models "Backend/models"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

func CMD(comando string) (bytes.Buffer, string, error) { //(entrada,salida)

	var salida bytes.Buffer
	var errors bytes.Buffer
	cmd := exec.Command("bash", "-c", comando)
	cmd.Stdout = &salida
	cmd.Stderr = &errors
	err := cmd.Run()
	return salida, errors.String(), err
}

func RequestPrincipal() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		salida, _, verificar := CMD("cat /proc/cpu_grupo18")
		fmt.Println(salida)

		if verificar != nil {
			log.Printf("error: %v\n", verificar)
		} else {
			var dataJson Models.CPUDATAJSON
			json.Unmarshal(salida.Bytes(), &dataJson) //json a objeto
			json.NewEncoder(rw).Encode(dataJson)
		}
	}
}

func RequestKill() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/Kill" {
			http.NotFound(rw, r)
			return
		}
		if  r.Method =="POST" {
			id := r.URL.Query().Get("pid")
			id = strings.TrimSuffix(id, "/")
      fmt.Println(id)
			  _, _, verificar := CMD("sudo kill -9 " + id)

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

func RequestCPU() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ejecutar := "ps -eo pcpu | sort -k 1 -r | head -50"
		salida, _, verificar := CMD(ejecutar)

		if verificar != nil {
			log.Printf("error: %v\n", verificar)
		} else {
			result := strings.Split(salida.String(), "\n")
			// fmt.Println(result)
			var theArray [50]string
			entro := false
			cont := 0
			for _, item := range result {

				if !entro {
					entro = true
				} else {
					theArray[cont] += strings.TrimSpace(item)
					cont++
				}
			}
			// fmt.Println(errout)
			json.NewEncoder(rw).Encode(theArray)
		}
	}
}
func RequestMemory() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		salida, _, verificar := CMD("cat /proc/mem_grupo18")

		if verificar != nil {
			log.Printf("error: %v\n", verificar)
		} else {

			var dataJson Models.DATAJSONMEMORY
			json.Unmarshal(salida.Bytes(), &dataJson) //json a objeto
			ejecutar := "free -m"
			salidaMsj, _, verificar := CMD(ejecutar) //para cache

			if verificar != nil {
				log.Printf("error: %v\n", verificar)
			} else {
				result := strings.Split(salidaMsj.String(), "\n")
				memoria := strings.ReplaceAll(result[1], " ", ",") //fila memoria
				cacheBusqueda := strings.Split(memoria, ",")       //buscar valor cache
				cont := 0
				cachestr := "0"
				for _, item := range cacheBusqueda {
					if item != "" {
						cont++
						if cont == 6 {
							cachestr = item
						}
					}
				}

				valor, err := strconv.Atoi(cachestr)
				if err != nil {
					fmt.Println(err)
				} else {
					dataJson.CACHE = valor * 1000000
				}
			}
      rw.Header().Set("Content-Type", "application/json")
			json.NewEncoder(rw).Encode(dataJson)
		}
	}
}
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API GO!\n"))
}
