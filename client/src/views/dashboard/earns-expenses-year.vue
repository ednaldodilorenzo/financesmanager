<template>
  <div class="card chart-card border-0 shadow-sm h-100">
    <div class="card-header bg-white border-bottom p-3">
      <h2 class="h6 fw-bold mb-1">Evolução no ano</h2><span class="small text-body-secondary">Receitas, despesas e
        investimentos até {{ selectedMonth }}</span>
    </div>
    <div class="card-body chart-body">
      <div class="chart-container">
        <LineChart v-if="dataList.length" :chart-data="chartData" :options="chartOptions" />
        <div v-else class="h-100 d-flex flex-column align-items-center justify-content-center text-body-secondary"><i
            class="bi bi-calendar3 fs-2"></i><span class="small mt-2">Sem dados para este ano</span></div>
      </div>
    </div>
  </div>
</template>
<script setup>
import { computed } from "vue";
import { LineChart } from "vue-chart-3";
import { currencyBRL } from "@/components/filters/currency.filter";
const props = defineProps({ dataList: { type: Array, default: () => [] }, date: { type: Date, default: () => new Date() } });
const allMonths = ["Jan", "Fev", "Mar", "Abr", "Mai", "Jun", "Jul", "Ago", "Set", "Out", "Nov", "Dez"];
const monthCount = computed(() => props.date.getMonth() + 1);
const labels = computed(() => allMonths.slice(0, monthCount.value));
const selectedMonth = computed(() => new Intl.DateTimeFormat("pt-BR", { month: "long" }).format(props.date));
function typeValues(type) { return props.dataList.filter((item) => item.type === type).slice(0, monthCount.value).map((item) => Math.abs(Number(item.total) || 0) / 100); }
const chartData = computed(() => ({
  labels: labels.value, datasets: [
    { label: "Receitas", data: typeValues("R"), borderColor: "#10b981", backgroundColor: "rgba(16,185,129,.08)", tension: .35, fill: false },
    { label: "Despesas", data: typeValues("D"), borderColor: "#ef4444", backgroundColor: "rgba(239,68,68,.08)", tension: .35, fill: false },
    { label: "Investimentos", data: typeValues("I"), borderColor: "#4f46e5", backgroundColor: "rgba(79,70,229,.08)", tension: .35, fill: false },
  ]
}));
const chartOptions = { responsive: true, maintainAspectRatio: false, interaction: { mode: "index", intersect: false }, plugins: { legend: { position: "bottom", labels: { usePointStyle: true } }, tooltip: { callbacks: { label: (context) => ` ${context.dataset.label}: ${currencyBRL(context.raw)}` } } }, scales: { x: { grid: { display: false } }, y: { beginAtZero: true, ticks: { callback: (value) => currencyBRL(value) } } } };
</script>
<style scoped>
.chart-card {
  border-radius: .85rem
}

.chart-card .card-header {
  border-radius: .85rem .85rem 0 0
}

.chart-body {
  height: 20rem;
  position: relative
}

.chart-container {
  position: relative;
  width: 100%;
  height: 260px;
  max-height: 260px;
}

.chart-container :deep(canvas) {
  width: 100% !important;
  height: 100% !important;
  max-width: 100% !important;
  max-height: 260px !important;
}
</style>
