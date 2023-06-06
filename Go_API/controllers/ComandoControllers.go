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

    rw.Header().Set("Content-Type", "application/json")
		salida, _, verificar := CMD("cat /proc/cpu_grupo18")
		
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
		if  r.Method =="GET" {
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

func RequestMemory() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		salida, _, verificar := CMD("cat /proc/mem_grupo18")

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
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API GO!\n"))
}
