package Models

type PROCESOPADRE struct {
	ID            string        `json:"pid"`
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
  CPU_USAGE     string `json:cpu_usage`
	DATA          []PROCESOPADRE
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
