<template>
  <div class="card">
    <div class="card-body">
      <h4 class="card-title">Receitas x Despesas Dia a Dia</h4>
      <LineChart
        :chart-data="chartLineEarnsExpensesEvolutionData"
        :options="chartLineOptions"
      ></LineChart>
    </div>
  </div>
</template>
<script setup>
import { computed } from 'vue';
import { LineChart } from "vue-chart-3";
import { getDaysListPerMonth } from "@/utils/date";

const props = defineProps({
  transactionsList: {
    type: Array,
    default: () => [],
  },
  date: {
    type: Date,
    default: () => new Date(),
  },
});

const chartLineEarnsExpensesEvolutionData = computed(() => ({
  labels: getDaysListPerMonth(props.date),
  datasets: [
    {
      label: "Despesa",
      data: evolutionDayByDay.value.map((value) =>
        Math.abs(value.expensesSum / 100)
      ),
      borderColor: "#FF0000",
      backgroundColor: "rgba(213, 63, 21, 0.2)",
      fill: true,
      tension: 0.4, // Smooth curve
    },
    {
      label: "Receita",
      data: evolutionDayByDay.value.map((value) => value.earnsSum / 100),
      borderColor: "#42A5F5",
      backgroundColor: "rgba(66, 165, 245, 0.2)",
      fill: true,
      tension: 0.4, // Smooth curve
    },
  ],
}));

const evolutionDayByDay = computed(() => {
  const orderedSummary = props.transactionsList.reduce(
    (previous, current) => {
      const day = new Date(current.paymentDate).getUTCDate();
      const item = previous.find((value) => value.day === day);
      if (current.category.type === "D") {
        item["expensesSum"] += current.value;
      } else if (current.category.type === "R") {
        item["earnsSum"] += current.value;
      }

      return previous;
    },
    getDaysListPerMonth(props.date)
      .filter((value) => value <= props.date.getUTCDate())
      .map((value) => ({
        day: value,
        expensesSum: 0.0,
        earnsSum: 0.0,
      }))
  );

  return orderedSummary.reduce((previous, current) => {
    if (previous.length > 0) {
      const lastInserted = previous[previous.length - 1];
      previous.push({
        day: current.day,
        expensesSum: lastInserted.expensesSum + current.expensesSum,
        earnsSum: lastInserted.earnsSum + current.earnsSum,
      });
    } else {
      previous.push({
        day: current.day,
        expensesSum: current.expensesSum,
        earnsSum: current.earnsSum,
      });
    }
    return previous;
  }, []);
});

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
