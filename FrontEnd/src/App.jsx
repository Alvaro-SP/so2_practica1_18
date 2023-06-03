import "./App.css";
import { data } from "./mocks/prueba.json"
import { ProcessTable } from "./components/Table";
import { useEffect, useState } from "react";

const API = "http://localhost:4000/"

function App() {
  const [procesos,setProcesos] = useState({cpu_usage:0,data})
  useEffect(()=>{
    fetch(`${API}Procesos`)
    .then(res => res.json())
    .then(data => setProcesos(data))
    .catch(err => console.log(err.message))
  },[])
  return (
    <>
      <header>
        <h1>Práctica Sistemas Operativos 2</h1>
        <h2>Grupo 18</h2>
      </header>
      <main>
        <h3>Procesos</h3>
        <h3>Tabla de procesos</h3>
        <ProcessTable data={procesos.data}/>
      </main>
    </>
  );
}

export default App;
