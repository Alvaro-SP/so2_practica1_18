import "./Memoria.css";
export function Memoria({ data }) {
  return (
    <>
      <div className="memoria">
        <Field title={"Total"} text={data.total/(1024*1024)} />
        <Field title={"Libre"} text={data.libre/(1024*1024)} />
        <Field title={"Porcentaje"} text={100-(data.libre/data.total)} />
      </div>
      <div className="memoria">
        <Field title={"Total"} text={data.total} />
        <Field title={"Libre"} text={data.libre} />
        <Field title={"Buffer"} text={data.buffer} />
        <Field title={"Cache"} text={data.cache} />
        <Field title={"Unidad"} text={data.unidad} />
      </div>
    </>
  );
}
export function Field({ title, text }) {
  return (
    <div>
      <h4>{title}</h4>
      <p>{isNaN(text)?text:text.toFixed(4)}</p>
    </div>
  );
}
