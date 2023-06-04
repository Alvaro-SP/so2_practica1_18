package Models

type PROCESOPADRE struct {
	ID            string        `json:"id"`
	ESTADO        string        `json:"estado"`
	NOMBRE        string        `json:"nombre"`
	PID           int           `json:"pid"`
	PROCESOSHIJOS []PROCESOHIJO `json:"procesos_hijo"`
	RAM           int           `json:"ram"`
	USUARIO       string        `json:"usuario"`
}

type PROCESOHIJO struct {
	ESTADO    string `json:"estado"`
	NOMBRE    string `json:"nombre"`
	PID       int    `json:"pid"`
	RAM       int    `json:"ram"`
	USUARIO   string `json:"usuario"`
}

type CPUDATAJSON struct {
  	CPU_USAGE     int `json:cpu_usage`
	DATA          []PROCESOPADRE
	Ejecucion 	int		"json:ejecucion"
    Zombie 		int		"json:zombie"
    Detenido 	int		"json:detenido"
    Suspendid 	int		"json:suspendido"
    Totales 	int		"json:totales"
}

type DATAJSONCPU struct {
	PORCENTAJE_PROCESO []string `json:"porcentaje_proceso"`
}

type DATAJSONMEMORY struct {
	MEMORIA_TOTAL int    `json:"memoria_total"`
	MEMORIA_LIBRE int    `json:"memoria_libre"`
	BUFFER        int    `json:"buffer"`
	CACHE         int    `json:"cache"`
	MEM_UNIT      int    `json:"mem_unit"`
	PORCENTAJE    int    `json:"porcentaje"`
}
