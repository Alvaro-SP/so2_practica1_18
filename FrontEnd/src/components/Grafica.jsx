import { LineChart } from "@tremor/react";

export function Grafica({datos}) {
  const dataFormatter = ( number) =>
  `${Intl.NumberFormat("us").format(number).toString()}%`;
  return (
    <LineChart
      className="mt-6"
      data={datos}
      index="ejex"
      categories={["porcentaje"]}
      colors={["emerald"]}
      yAxisWidth={40}
      autoMinValue
      showXAxis={false}
      valueFormatter={dataFormatter}
    />
  );
}
