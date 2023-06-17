import React, { useState } from "react";
import "./Table.css";
import { useEffect } from "react";
import maps from "../mocks/maps.json";

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
    fetch(`${API}Kill?pid=${pid}`)
      .then((res) => console.log(res))
      .catch((err) => console.log(err));
  };
  return (
    <>
      {showModal && <ModalRam pid={pid} cerrarModal={showRam} />}
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
export function ModalRam({ pid, cerrarModal }) {
  const [asignaciones, setAsignaciones] = useState(maps);
  useEffect(() => {
    // Fetch para obtener maps
    fetch(`${API}maps?pid=${pid}`)
      .then((res) => res.json())
      .then((maps) => {
        console.log(maps);
        setAsignaciones(maps);
      })
      .catch((err) => console.log(err));
  }, []);
  return (
    <section className="modal-overlay">
      <div className="modal">
        <h2>{pid} - Asignación de memoria</h2>
        <section style={{ height: "50vh", overflowY: "scroll" }}>
          <table>
            <thead>
              <th>Direccion</th>
              <th>Tamaño</th>
              <th>Permisos</th>
              <th>Dispositivo</th>
              <th>Archivo</th>
            </thead>
            <tbody>
              {asignaciones.map((value) => (
                <tr className={"childrow"} key={value.direccion}>
                  <td>{value.direccion}</td>
                  <td>{value.tamanio}</td>
                  <td>{value.permisos}</td>
                  <td>{value.dispositivo}</td>
                  <td>{value.archivo}</td>
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
