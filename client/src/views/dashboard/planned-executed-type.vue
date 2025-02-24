<template>
  <div class="card">
    <div class="card-body">
      <h4 class="card-title">Planejado x Executado por Tipo</h4>
      <BarChart :chart-data="chartData"> </BarChart>
    </div>
  </div>
</template>
<script setup>
import { computed } from "vue";
import { BarChart } from "vue-chart-3";

const props = defineProps({
  dataList: {
    type: Array,
    default: () => [],
  },
});

const chartData = computed(() => ({
  labels: ["Receitas", "Despesas", "Investimentos"],
  datasets: [
    {
      label: "Planejado",
      data: [
        chartValues.value.earns.planned / 100,
        chartValues.value.expenses.planned / 100,
        chartValues.value.investments.planned / 100,
      ],
      borderColor: "#06c2f4",
      backgroundColor: "#06c2f4",
      fill: true,
      tension: 0.4,
    },
    {
      label: "Executado",
      data: [
        chartValues.value.earns.executed / 100,
        Math.abs(chartValues.value.expenses.executed / 100),
        Math.abs(chartValues.value.investments.executed / 100),
      ],
      borderColor: "#3daa02",
      backgroundColor: "#3daa02",
      fill: true,
      tension: 0.4,
    },
  ],
}));

const chartValues = computed(() =>
  props.dataList.reduce(
    (previous, current) => {
      if (current.type === "D") {
        previous.expenses.planned += current.planned / 12;
        previous.expenses.executed += current.total;
      } else if (current.type === "R") {
        previous.earns.planned += current.planned / 12;
        previous.earns.executed += current.total;
      } else {
        previous.investments.planned += current.planned / 12;
        previous.investments.executed += current.total;
      }
      return previous;
    },
    {
      earns: {
        planned: 0.0,
        executed: 0.0,
      },
      expenses: {
        planned: 0.0,
        executed: 0.0,
      },
      investments: {
        planned: 0.0,
        executed: 0.0,
      },
    }
  )
);
</script>
