import "./Memoria.css";
export function Memoria({ data }) {
  return (
    <>
      <div className="memoria">
        <Field title={"Total"} text={(data.total/(1024*1024)).toFixed(2)} />
        <Field title={"Libre"} text={(data.libre/(1024*1024)).toFixed(2)} />
        <Field title={"Porcentaje"} text={data.porcentaje} />
      </div>
    </>
  );
}
export function Field({ title, text }) {
  return (
    <div>
      <h4>{title}</h4>
      <p>{text}</p>
    </div>
  );
}
