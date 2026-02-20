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
      data: props.dataList.filter((item) => item.type === "D").map((item) => Math.abs(item.total / 100)),
      borderColor: "#FF0000",
      backgroundColor: "rgba(213, 63, 21, 0.2)",
      fill: true,
      tension: 0.4, // Smooth curve
    },
    {
      label: "Receita",
      data: props.dataList.filter((item) => item.type === "R").map((item) => item.total / 100),
      borderColor: "#42A5F5",
      backgroundColor: "rgba(66, 165, 245, 0.2)",
      fill: true,
      tension: 0.4, // Smooth curve
    },
    {
      label: "Investimentos",
      data: props.dataList.filter((item) => item.type === "I").map((item) => -item.total / 100),
      borderColor: "#4CAF50", // Green 500
      backgroundColor: "rgba(76, 175, 80, 0.2)", // Green 500 with 20% opacity
      fill: true,
      tension: 0.4, // Smooth curve
    },
  ],
}));

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
