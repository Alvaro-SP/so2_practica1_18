import "./App.css";
import { DATA } from "./mocks/prueba.json";
import { memo } from "./mocks/memoria.json";
import { ProcessTable } from "./components/Table";
import { useEffect, useState } from "react";
import { Memoria } from "./components/Memoria";
import { Grafica } from "./components/Grafica";
import { Procesos } from "./components/Procesos";

const API = import.meta.env.VITE_API 

function App() {
  const [procesos, setProcesos] = useState({
    CPU_USAGE: 0,
    DATA,
    Detenido: 28,
    Ejecucion: 10,
    Totales: 100,
    Suspendid: 16,
    Zombie: 0,
  });
  const [memoria, setMemoria] = useState(memo);
  const [historialMemoria, setHistorialMemoria] = useState([]);
  useEffect(() => {
    const xd = [];
    const id = setInterval(() => {
      fetch(`${API}Memoria`)
        .then((res) => res.json())
        .then((data) => {
          const mapMemoria = {
            total: data.memoria_total,
            libre: data.memoria_libre,
            buffer: data.buffer,
            porcentaje: data.porcentaje / 100,
            unidad: data.mem_unit,
            ejex: xd.length,
          };
          xd.push(mapMemoria);
          setMemoria(mapMemoria);
          setHistorialMemoria([...xd]);
        })
        .catch((err) => console.log(err.message));
      fetch(`${API}Principal`)
      .then(res => res.json())
      .then(data => setProcesos(data))
      .catch(er => console.log(er))
    }, 3000);
    return () => {
      clearInterval(id);
    };
  }, []);
  return (
    <>
      <header>
        <h1>Práctica Sistemas Operativos 2</h1>
        <h2>Grupo 18</h2>
      </header>
      <main>
        <h1>Memoria RAM (MB)</h1>
        <Memoria data={memoria} />
        <h3>Consumo de memoria</h3>
        <Grafica memoria={memoria} datos={historialMemoria} />
        <h1>Procesos</h1>
        <Procesos
          total={procesos.Totales}
          exe={procesos.Ejecucion}
          suspendidos={procesos.Suspendid}
          detenidos={procesos.Detenido}
          zombies={procesos.Zombie}
        />
        <h3>Tabla de procesos</h3>
        <ProcessTable data={procesos.DATA} />
      </main>
    </>
  );
}

export default App;
