import { LineChart } from "@tremor/react";

export function Grafica({ datos }) {
  const dataFormatter = (number) => `${number.toFixed(2).toString()}MB`;
  return (
    <LineChart
      className="mt-6"
      data={datos}
      index="ejex"
      categories={["consumo"]}
      colors={["emerald"]}
      yAxisWidth={60}
      autoMinValue
      showXAxis={false}
      valueFormatter={dataFormatter}
    />
  );
}
