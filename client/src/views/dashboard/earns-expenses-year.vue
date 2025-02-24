<template>
  <div class="card">
    <div class="card-body">
      <h4 class="card-title">Receitas x Despesas no Ano</h4>
      <LineChart
        :chart-data="chartLineData"
        :options="chartLineOptions"
      ></LineChart>
    </div>
  </div>
</template>
<script setup>
import { ref, computed } from "vue";
import { LineChart } from "vue-chart-3";
import { getMonthsListUntilDate } from "@/utils/date";

const props = defineProps({
  dataList: {
    type: Array,
    default: () => [],
  },
  date: {
    type: Date,
    default: () => new Date(),
  },
});
const monthsList = ref(getMonthsListUntilDate(props.date));

// Computed chart data
const chartLineData = computed(() => ({
  labels: monthsList.value,
  datasets: [
    {
      label: "Despesa",
      data: chartValues.value.map((item) => Math.abs(item.expenses / 100)),
      borderColor: "#FF0000",
      backgroundColor: "rgba(213, 63, 21, 0.2)",
      fill: true,
      tension: 0.4, // Smooth curve
    },
    {
      label: "Receita",
      data: chartValues.value.map((item) => item.earns / 100),
      borderColor: "#42A5F5",
      backgroundColor: "rgba(66, 165, 245, 0.2)",
      fill: true,
      tension: 0.4, // Smooth curve
    },
  ],
}));

const chartValues = computed(() =>
  props.dataList.reduce(
    (previous, current) => {
      if (current.type === "D") {
        previous.forEach((month) => {
          month.expenses += current.accumulated / monthsList.value.length;
        });
      } else if (current.type === "R") {
        previous.forEach((month) => {
          month.earns += current.accumulated / monthsList.value.length;
        });
      }
      return previous;
    },
    monthsList.value.map((value) => ({ earns: 0.0, expenses: 0.0 }))
  )
);

// Chart options
const chartLineOptions = {
  responsive: true,
  maintainAspectRatio: false,
  scales: {
    y: {
      beginAtZero: true,
    },
  },
};
</script>
