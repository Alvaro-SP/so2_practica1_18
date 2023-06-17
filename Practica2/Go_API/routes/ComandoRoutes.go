package routes

import (
	"Backend/controllers"
	"github.com/gorilla/mux"
)

func ComandoRoute(router *mux.Router) {
	//Rutas que serán consumidas por el FrontEnd
	router.HandleFunc("/", Controllers.IndexHandler).Methods("GET") //Retorna un mensaje de prueba
	router.HandleFunc("/Principal", Controllers.RequestPrincipal()) //Obtiene el árbol de procesos e información de cpu
	router.HandleFunc("/Kill", Controllers.RequestKill())           //Recibe un parámetro pid para ejecutar el comando kill
	router.HandleFunc("/Memoria", Controllers.RequestMemory())      //Retorna un objeto con información actual de la memoria
	router.HandleFunc("/maps", Controllers.RequestMaps()) 			//Obtenemos la administración de memoria de un proceso, recibe un parámetro (pid)
}
