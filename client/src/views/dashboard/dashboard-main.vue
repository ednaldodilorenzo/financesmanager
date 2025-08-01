<template>
  <div class="mt-3">
    <div class="d-flex justify-content-between">
      <h2 class="fs-4">Dashboard</h2>
      <nav style="--bs-breadcrumb-divider: '>'" aria-label="breadcrumb">
        <ol class="breadcrumb">
          <li class="breadcrumb-item"><a href="#">Home</a></li>
          <li class="breadcrumb-item active"><a href="#">Dashboard</a></li>
        </ol>
      </nav>
    </div>
  </div>
  <div class="card mb-3">
    <div class="card-body p-2">
      <div class="d-flex justify-content-center my-3">
        <h3>{{ getExtenseMonth(currentDate) }}</h3>
      </div>
    </div>
  </div>
  <div class="row g-3">
    <div class="col-4">
      <div class="card">
        <div class="card-body">
          <h4 class="card-title">Valores por Tipo</h4>
          <PieChart :chart-data="chartData" :options="chartOptions"></PieChart>
        </div>
      </div>
    </div>
    <div class="col-4">
      <div class="card">
        <div class="card-body">
          <h4 class="card-title">Despesas por Categoria</h4>
          <PieChart
            :chart-data="chartExpensesData"
            :options="chartOptions"
          ></PieChart>
        </div>
      </div>
    </div>
    <div class="col-4">
      <div class="card">
        <div class="card-body">
          <h4 class="card-title">Receitas por Categoria</h4>
          <PieChart
            :chart-data="chartEarnsData"
            :options="chartOptions"
          ></PieChart>
        </div>
      </div>
    </div>
    <div class="col-4">
      <EarnsExpensesYear
        :date="currentDate"
        :data-list="plannedList"
      ></EarnsExpensesYear>
    </div>
    <div class="col-4">
      <day-by-day-graph
        :transactionsList="transactionsList"
        :date="currentDate"
      ></day-by-day-graph>
    </div>
    <div class="col-4">
      <PlannedExecutedType :dataList="plannedList"></PlannedExecutedType>
    </div>
  </div>
</template>
<script setup>
import { ref, computed } from "vue";
import { PieChart } from "vue-chart-3";
import { Chart, registerables } from "chart.js";
import transactionService from "@/views/transaction/transaction.service";
import { useLoadingScreen } from "@/components/loading/useLoadingScreen";
import { getExtenseMonth } from "@/utils/date";
import DayByDayGraph from "./daybyday-graph.vue";
import planningService from "../planning/planning.service";
import PlannedExecutedType from "./planned-executed-type.vue";
import EarnsExpensesYear from "./earns-expenses-year.vue";
import { useRouter } from "vue-router";

// Register required components
Chart.register(...registerables);

const loading = useLoadingScreen();
const currentDate = new Date();
const router = useRouter();

const transactionsList = ref([]);
const plannedList = ref([]);

function loadInitalData() {
  loading.show();
  const params = {
    month: currentDate.getMonth() + 1,
    year: currentDate.getFullYear(),
  };

  Promise.allSettled([
    transactionService.findAll(params),
    planningService.findAll(params),
  ])
    .then((results) => {
      const [transactionsResult, planningResults] = results;
      transactionsList.value = transactionsResult.value.data;
      plannedList.value = planningResults.value.data;
    })
    .catch((err) => {
      router.push({ name: "denied" });
    })
    .finally(() => {
      loading.hide();
    });
}

const chartOptions = {
  plugins: {
    legend: {
      display: false, // Remove legend
    },
    tooltip: {
      enabled: true, // Keep tooltips
    },
  },
};

loadInitalData();

const valuesByTypeData = computed(() =>
  transactionsList.value.reduce(
    (previous, current) => {
      if (current.category.type === "D") {
        const result = [...previous];
        result[1] = result[1] + current.value / 100;
        return result;
      }
      if (current.category.type === "R") {
        const result = [...previous];
        result[0] = result[0] + current.value / 100;
        return result;
      } else {
        const result = [...previous];
        result[2] = result[2] + Math.abs(current.value / 100);
        return result;
      }
    },
    [0.0, 0.0, 0.0]
  )
);

const valuesByCategory = computed(() =>
  transactionsList.value.reduce((previous, current) => {
    const item = previous.find(
      (value) => value.category === current.category.name
    );
    if (item) {
      item.value += current.value / 100;
    } else {
      previous.push({
        category: current.category.name,
        value: current.value / 100,
        type: current.category.type,
      });
    }

    return previous;
  }, [])
);

const chartData = computed(() => ({
  labels: ["Receitas", "Despesas", "Investimentos"],
  datasets: [
    {
      label: "Sales",
      data: valuesByTypeData.value,
      backgroundColor: ["#f87979", "#a1d4ef", "#85e85e"],
    },
  ],
}));

const getRandomColor = () => {
  return `#${Math.floor(Math.random() * 0xffffff)
    .toString(16)
    .padStart(6, "0")}`;
};

const chartExpensesData = computed(() => ({
  labels: valuesByCategory.value
    .filter((item) => item.type === "D")
    .map((item) => item.category),
  datasets: [
    {
      label: "Despesas",
      data: valuesByCategory.value
        .filter((item) => item.type === "D")
        .map((item) => item.value),
      backgroundColor: valuesByCategory.value
        .filter((item) => item.type === "D")
        .map(() => getRandomColor()),
    },
  ],
}));

const chartEarnsData = computed(() => ({
  labels: valuesByCategory.value
    .filter((item) => item.type === "R")
    .map((item) => item.category),
  datasets: [
    {
      label: "Sales",
      data: valuesByCategory.value
        .filter((item) => item.type === "R")
        .map((item) => item.value),
      backgroundColor: valuesByCategory.value
        .filter((item) => item.type === "R")
        .map(() => getRandomColor()),
    },
  ],
}));
</script>
