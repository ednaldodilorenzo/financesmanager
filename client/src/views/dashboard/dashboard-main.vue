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
  <div class="row">
    <div class="col-4">
      <div class="card">
        <div class="card-body">
          <h4 class="card-title">Valores por tipo</h4>
          <PieChart :chart-data="chartData"></PieChart>
        </div>
      </div>
    </div>
    <div class="col-4">
      <div class="card">
        <div class="card-body">
          <h4 class="card-title">Despesas por Categoria</h4>
          <PieChart :chart-data="chartExpensesData"></PieChart>
        </div>
      </div>
    </div>
    <div class="col-4">
      <div class="card">
        <div class="card-body">
          <h4 class="card-title">Receitas por Categoria</h4>
          <PieChart :chart-data="chartEarnsData"></PieChart>
        </div>
      </div>
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

// Register required components
Chart.register(...registerables);

const loading = useLoadingScreen();
const currentDate = new Date();

const transactionsList = ref([]);

function loadInitalData() {
  loading.show();
  const params = {
    month: currentDate.getMonth() + 1,
    year: currentDate.getFullYear(),
  };

  transactionService
    .findAll(params)
    .then((resp) => {
      transactionsList.value = resp.items;
    })
    .catch((err) => {
      router.push({ name: "denied" });
    })
    .finally(() => {
      loading.hide();
    });
}

loadInitalData();

const valuesByTypeData = computed(() =>
  transactionsList.value.reduce(
    (previous, current) => {
      if (current.category.type === "D") {
        const result = [...previous];
        result[1] = result[1] + Math.abs(current.value / 100);
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
      item.value += Math.abs(current.value / 100);
    } else {
      previous.push({
        category: current.category.name,
        value: Math.abs(current.value / 100),
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
  const letters = "0123456789ABCDEF";
  let color = "#";
  for (let i = 0; i < 6; i++) {
    color += letters[Math.floor(Math.random() * 16)];
  }
  return color;
};

const chartExpensesData = computed(() => ({
  labels: valuesByCategory.value
    .filter((item) => item.type === "D")
    .map((item) => item.category),
  datasets: [
    {
      label: "Sales",
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
