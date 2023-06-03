import "./App.css";
import { data } from "./mocks/prueba.json";
import { memo } from "./mocks/memoria.json";
import { ProcessTable } from "./components/Table";
import { useEffect, useState } from "react";
import { Memoria } from "./components/Memoria";

const API = "http://localhost:4000/";

function App() {
  const [procesos, setProcesos] = useState({ cpu_usage: 0, data });
  const [memoria, setMemoria] = useState(memo);
  useEffect(() => {
    console.log(memoria);
    fetch(`${API}Memoria`)
      .then((res) => res.json())
      .then((data) => {
        const mapMemoria = {
          total: data.memoria_total,
          libre: data.memoria_libre,
          buffer: data.buffer,
          cache: data.cache,
          unidad: data.mem_unit,
        };
        setMemoria(mapMemoria);
      })
      .catch((err) => console.log(err.message));
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
        <h3>Procesos</h3>
        <h3>Tabla de procesos</h3>
        <ProcessTable data={procesos.data} />
      </main>
    </>
  );
}

export default App;
