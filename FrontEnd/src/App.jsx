import "./App.css";
import { data } from "./mocks/prueba.json";
import { memo } from "./mocks/memoria.json";
import { ProcessTable } from "./components/Table";
import { useEffect, useState } from "react";
import { Memoria } from "./components/Memoria";
import { Grafica } from "./components/Grafica";

const API = "http://localhost:4000/";

function App() {
  const [procesos, setProcesos] = useState({ cpu_usage: 0, data });
  const [memoria, setMemoria] = useState(memo);
  const [historialMemoria, setHistorialMemoria] = useState([]);
  useEffect(() => {
      const xd = []
    const id = setInterval(() => {
      fetch(`${API}Memoria`)
        .then((res) => res.json())
        .then((data) => {
          const mapMemoria = {
            total: data.memoria_total,
            libre: data.memoria_libre,
            buffer: data.buffer,
            porcentaje: data.porcentaje/100,
            unidad: data.mem_unit,
            ejex: xd.length,
          };
          xd.push(mapMemoria)
          setMemoria(mapMemoria)
          setHistorialMemoria([...xd])
        })
        .catch((err) => console.log(err.message));
    }, 3000);
    return ()=>{
      clearInterval(id)
    }
  }, []);
  return (
    <>
      <header>
        <h1>Práctica Sistemas Operativos 2</h1>
        <h2>Grupo 18</h2>
      </header>
      <main>
        <h3>Memoria RAM (MB)</h3>
        <Memoria data={memoria} />
        <h3>Consumo de memoria</h3>
        <Grafica memoria={memoria} datos={historialMemoria} />
        <h2>Procesos</h2>
        <h3>Tabla de procesos</h3>
        <ProcessTable data={procesos.data} />
      </main>
    </>
  );
}

export default App;
