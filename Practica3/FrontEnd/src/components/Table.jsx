import React, { useState } from "react";
import "./Table.css";
import { useEffect } from "react";
import maps from "../mocks/maps.json";
import axios from "axios";
import { Field } from "./Memoria";
const API = import.meta.env.VITE_API;
export function ProcessTable({ data }) {
  return (
    <section style={{ height: "70vh", overflowY: "scroll" }}>
      <table>
        <thead>
          <tr>
            <th>PID</th>
            <th>Nombre</th>
            <th>Usuario</th>
            <th>Estado</th>
            <th>%RAM</th>
            <th colSpan={2}>Acción</th>
          </tr>
        </thead>
        <tbody>
          {data.map((value) => {
            return <ParentRow key={value.id} {...value} />;
          })}
        </tbody>
      </table>
    </section>
  );
}
export function ParentRow(
  { pid, nombre, usuario, estado, ram, procesoshijos },
) {
  const [isExpanded, setIsExpanded] = useState(false);
  const [showModal, setShowModal] = useState(false);

  const handleClick = () => {
    setIsExpanded(!isExpanded);
  };
  const showRam = (e) => {
    e.stopPropagation();
    setShowModal(!showModal);
  };

  const sendKill = (e) => {
    e.stopPropagation();
    console.log(pid);
    axios.get(`${API}Kill?pid=${pid}`)
      .then((res) => console.log(res))
      .catch((err) => console.log(err));
  };
  return (
    <>
      {showModal && (
        <ModalRam
          pid={pid}
          nombre={nombre}
          cerrarModal={showRam}
        />
      )}
      <tr
        className={isExpanded ? "expanded" : "no-expanded"}
        onClick={handleClick}
      >
        <td>{pid}</td>
        <td>{nombre}</td>
        <td>{usuario}</td>
        <td>{estado}</td>
        <td>{ram / 1000} %</td>
        <td>
          <button onClick={showRam} className={"btn-kill"}>Ver RAM</button>
        </td>
        <td>
          <button onClick={sendKill} className={"btn-kill"}>Kill</button>
        </td>
      </tr>
      {isExpanded && procesoshijos &&
        procesoshijos.map((value) => <ChildRow {...value} key={value.pid} />)}
    </>
  );
}
export function ChildRow({ pid, nombre, usuario, estado, ram }) {
  return (
    <tr className={"childrow"}>
      <td>{pid}</td>
      <td>{nombre}</td>
      <td>{usuario}</td>
      <td>{estado}</td>
      <td>{ram / 1000}%</td>
    </tr>
  );
}
export function ModalRam({ pid, cerrarModal, nombre }) {
  const [asignaciones, setAsignaciones] = useState([]);
  const [memoria, setMemoria] = useState({ mr: 0, mv: 1 });
  useEffect(() => {
    // Post para obtener maps
    axios.get(`${API}maps?pid=${pid}`)
      .then((res) => res.data)
      .then((maps) => {
        if (maps == null) return;
        const { mapped, rss, size } = mapPermisos(maps.arr1, maps.arr2);
        setAsignaciones(mapped);
        setMemoria({ mv: size, mr: rss });
      })
      .catch((err) => console.log(err));
  }, []);
  const mapPermisos = (data, smaps) => {
    const rss = 0;
    const size = 0;
    const mapped = data.map((value, index) => {
      const listaPermisos = [];
      if (value.Permisos.includes("r")) listaPermisos.push("Lectura");
      if (value.Permisos.includes("w")) listaPermisos.push("Escritura");
      if (value.Permisos.includes("x")) listaPermisos.push("Ejecución");
      if (value.Permisos.includes("p")) listaPermisos.push("Privado");
      if (value.Permisos.includes("s")) listaPermisos.push("Compartido");
      rss += smaps[index].rss;
      size += smaps[index].size;
      return ({ ...value, Permisos: listaPermisos.join(",") });
    });
    return { mapped, rss, size };
  };
  const getInicio = () => {
    if (asignaciones.length == 0) return "-";
    const inicio = asignaciones[0]?.Direccion.split("-")[0];
    return inicio;
  };
  const getFin = () => {
    if (asignaciones.length == 0) return "-";
    const fin = asignaciones[asignaciones.length - 1]?.Direccion.split("-")[1];
    return fin;
  };
  return (
    <section className="modal-overlay">
      <div className="modal">
        <h2>{pid} {nombre} - Asignación de memoria</h2>
        <BarraMemoria
          mv={memoria.mv}
          mr={memoria.mr}
          inicio={getInicio()}
          fin={getFin()}
        />
        <section
          style={{ margin: "10px 0", height: "40vh", overflowY: "scroll" }}
        >
          <table>
            <thead>
              <th>Direccion</th>
              <th>Tamaño (Kb)</th>
              <th>Permisos</th>
              <th>Dispositivo</th>
              <th>Archivo</th>
              <th>RSS</th>
              <th>Size</th>
            </thead>
            <tbody>
              {mapPermisos(asignaciones).map((value, index) => (
                <tr className={"childrow"} key={index}>
                  <td>{value.Direccion}</td>
                  <td>{value.Tamanio / 1024}</td>
                  <td>{value.Permisos}</td>
                  <td>{value.Dispositivo}</td>
                  <td>{value.Archivo}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </section>
        <button onClick={cerrarModal}>Cerrar</button>
      </div>
    </section>
  );
}
function BarraMemoria({ mr = 0, mv = 1, inicio, fin }) {
  const getConsumo = () => mr / mv * 100;
  return (
    <section style={{ display: "flex", justifyContent: "space-evenly" }}>
      <div
        style={{
          display: "flex",
          flexDirection: "column",
          justifyContent: "center",
        }}
        className="cards"
      >
        <Field title={"Memoria residente"} text={mr + " MB"} />
        <Field title={"Memoria virtual"} text={mv + " MB"} />
        <Field title={"Consumo de Memoria"} text={getConsumo() + " %"} />
      </div>
      <div>
        <p>Mapa de memoria</p>
        <p className="border" style={{ marginTop: "5px" }}>Inicio: {inicio}</p>
        <div className="border" style={{ height: "30vh" }}>
          <div className="border-rss" style={{ height: `${getConsumo()}%` }}>
            RSS
          </div>
          VM
        </div>
        <p className="border">Fin: {fin}</p>
      </div>
    </section>
  );
}
