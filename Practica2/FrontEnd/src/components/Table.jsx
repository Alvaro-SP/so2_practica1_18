import React, { useState } from "react";
import "./Table.css";
import { useEffect } from "react";
import maps from "../mocks/maps.json";
import axios from "axios";
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
      {showModal && <ModalRam pid={pid} nombre={nombre} cerrarModal={showRam} />}
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
  useEffect(() => {
    // Post para obtener maps
    axios.get(`${API}maps?pid=${pid}`)
      .then((res) => res.data)
      .then((maps) => {
        setAsignaciones(maps??[]);
      })
      .catch((err) => console.log(err));
  }, []);
  const mapPermisos = (data) =>
    data.map((value) => {
      const listaPermisos = [];
      if (value.Permisos.includes("r")) listaPermisos.push("Lectura");
      if (value.Permisos.includes("w")) listaPermisos.push("Escritura");
      if (value.Permisos.includes("x")) listaPermisos.push("Ejecución");
      if (value.Permisos.includes("p")) listaPermisos.push("Privado");
      if (value.Permisos.includes("s")) listaPermisos.push("Compartido");
      return ({ ...value, Permisos: listaPermisos.join(",") });
    });
  return (
    <section className="modal-overlay">
      <div className="modal">
        <h2>{pid} {nombre} - Asignación de memoria</h2>
        <section
          style={{ margin: "10px 0", height: "50vh", overflowY: "scroll" }}
        >
          <table>
            <thead>
              <th>Direccion</th>
              <th>Tamaño (Kb)</th>
              <th>Permisos</th>
              <th>Dispositivo</th>
              <th>Archivo</th>
            </thead>
            <tbody>
              {mapPermisos(asignaciones).map((value, index) => (
                <tr className={"childrow"} key={index}>
                  <td>{value.Direccion}</td>
                  <td>{value.Tamanio}</td>
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
