package Models
// PROCESOPADRE representa un proceso principal en el sistema.
type PROCESOPADRE struct {
	ID            string        `json:"id"` // ID del proceso sirve para identificarlo
	ESTADO        string        `json:"estado"` // Estado del proceso sirve para tener el estado del proceso
	NOMBRE        string        `json:"nombre"` // Nombre del proceso sirve para tener el nombre del proceso
	PID           int           `json:"pid"` // PID del proceso sirve para tener el PID del proceso
	PROCESOSHIJOS []PROCESOHIJO `json:"procesoshijos"` // Procesos hijos del proceso sirve para tener los procesos hijos del proceso
	RAM           int           `json:"ram"` // RAM del proceso sirve para tener la RAM del proceso
	USUARIO       string        `json:"usuario"` // Usuario del proceso sirve para tener el usuario del proceso
}
// PROCESOHIJO representa un proceso hijo de un proceso principal.
type PROCESOHIJO struct {
	ESTADO    string `json:"estado"` // Estado del proceso sirve para guardar el estado del proceso hijo
	NOMBRE    string `json:"nombre"` // Nombre del proceso sirve para guardar el nombre del proceso hijo
	PID       int    `json:"pid"` // PID del proceso sirve para guardar el PID del proceso hijo
	RAM       int    `json:"ram"`	// RAM del proceso sirve para guardar la RAM del proceso hijo
	USUARIO   string `json:"usuario"` // Usuario del proceso sirve para guardar el usuario del proceso hijo
}
//  CPUDATAJSON representa la estructura de los datos de CPU.
type CPUDATAJSON struct {
  CPU_USAGE     int `json:cpu_usage` // CPU_USAGE sirve para guardar el uso de la CPU
  DATA          []PROCESOPADRE  // DATA sirve para guardar los datos de los procesos
	Ejecucion 	  int		"json:ejecucion" // Ejecucion sirve para guardar el numero de procesos en ejecucion
  Zombie 		    int		"json:zombie" // Zombie sirve para guardar el numero de procesos zombies
  Detenido 	    int		"json:detenido" // Detenido sirve para guardar el numero de procesos detenidos
  Suspendid 	  int		"json:suspendido" // Suspendido sirve para guardar el numero de procesos suspendidos
  Totales 	    int		"json:totales" // Totales sirve para guardar el numero de procesos totales
}
// DATAJSONCPU representa la estructura de los datos de CPU.
type DATAJSONCPU struct {
	PORCENTAJE_PROCESO []string `json:"porcentaje_proceso"` // PORCENTAJE_PROCESO sirve para guardar el porcentaje de uso de la CPU de cada proceso
}
// DATAJSONMEMORY representa la estructura de los datos de memoria.
type DATAJSONMEMORY struct {
	MEMORIA_TOTAL int    `json:"memoria_total"`  // MEMORIA_TOTAL sirve para guardar la memoria total del sistema
	MEMORIA_LIBRE int    `json:"memoria_libre"` // MEMORIA_LIBRE sirve para guardar la memoria libre del sistema
	BUFFER        int    `json:"buffer"` // BUFFER sirve para guardar el buffer del sistema
	MEM_UNIT      int    `json:"mem_unit"` // MEM_UNIT sirve para guardar la unidad de memoria del sistema
	PORCENTAJE    int    `json:"porcentaje"` // PORCENTAJE sirve para guardar el porcentaje de uso de la memoria del sistema
}
// MemoryMap representa la estructura de los datos de memoria.
type MemoryMap struct {
	Direccion  string // Direccion sirve para guardar la direccion de memoria
	Tamanio    uint64 // Tamanio sirve para guardar el tamanio de memoria
	Permisos   string  // Permisos sirve para guardar los permisos de memoria
	Dispositivo string // Dispositivo sirve para guardar el dispositivo de memoria
	Archivo    string	// Archivo sirve para guardar el archivo de memoria
}

// MemResVirtual representa la estructura de datos de Residente y virtual 
type MemResVirtual struct {
	Residente	uint64	// Memoria residente del proceso 
	Virtual		uint64	// Memoria virtual del proceso 
}
type Memory struct {
	Arr1          []MemoryMap
	Arr2          []MemResVirtual
}