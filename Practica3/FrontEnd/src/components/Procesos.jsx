import { Field } from "./Memoria";
import "./Memoria.css"

export function Procesos({total,exe,suspendidos,detenidos,zombies,uso}){
  return (
  <div className="memoria">
    <Field title={"Total procesos"} text={total}/>
    <Field title={"Ejecución"} text={exe}/>
    <Field title={"Suspendidos"} text={suspendidos}/>
    <Field title={"Detenidos"} text={detenidos}/>
    <Field title={"Zombies"} text={zombies}/>
    <Field title={"Porcentaje %"} text={uso}/>
  </div>
  )
}
