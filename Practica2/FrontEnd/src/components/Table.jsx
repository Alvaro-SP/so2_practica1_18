import React, { useState } from "react";
import "./Table.css";

const API = import.meta.env.VITE_API 
export function ProcessTable({ data }) {
  return (
    <table>
      <thead>
        <tr>
          <th>PID</th>
          <th>Nombre</th>
          <th>Usuario</th>
          <th>Estado</th>
          <th>%RAM</th>
          <th>Acción</th>
        </tr>
      </thead>
      <tbody>
        {data.map((value) => {
          return <ParentRow key={value.id} {...value} />;
        })}
      </tbody>
    </table>
  );
}
export function ParentRow(
  { pid, nombre, usuario, estado, ram, procesoshijos },
) {
  const [isExpanded, setIsExpanded] = useState(false);

  const handleClick = () => {
    setIsExpanded(!isExpanded);
  };

  const sendKill = ()=>{
    console.log(pid)
    fetch(`${API}Kill?pid=${pid}`)
    .then(res => console.log(res))
    .catch(err => console.log(err))
  }
  return (
    <>
      <tr
        className={isExpanded ? "expanded" : "no-expanded"}
        onClick={handleClick}
      >
        <td>{pid}</td>
        <td>{nombre}</td>
        <td>{usuario}</td>
        <td>{estado}</td>
        <td>{ram / 1000} %</td>
        <td><button onClick={sendKill} className={"btn-kill"}>Kill</button></td>
      </tr>
      {isExpanded && procesoshijos && procesoshijos.map((value) => <ChildRow {...value} />)}
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
