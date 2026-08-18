<template>
  <div class="card chart-card border-0 shadow-sm h-100">
    <div class="card-header bg-white border-bottom p-3">
      <h2 class="h6 fw-bold mb-1">Planejado x realizado</h2><span class="small text-body-secondary">Comparativo mensal
        por tipo</span>
    </div>
    <div class="card-body chart-body">
      <div class="chart-container">
        <BarChart v-if="dataList.length" :chart-data="chartData" :options="chartOptions" />
        <div v-else class="h-100 d-flex flex-column align-items-center justify-content-center text-body-secondary"><i
            class="bi bi-bar-chart fs-2"></i><span class="small mt-2">Sem planejamento para este mês</span></div>
      </div>
    </div>
  </div>
</template>
<script setup>
import { computed } from "vue";
import { BarChart } from "vue-chart-3";
import { currencyBRL } from "@/components/filters/currency.filter";
const props = defineProps({ dataList: { type: Array, default: () => [] } });
const totals = computed(() => props.dataList.reduce((sum, item) => {
  const key = item.type === "R" ? "earns" : item.type === "D" ? "expenses" : "investments";
  sum[key].planned += Math.abs(Number(item.planned) || 0) / 12 / 100;
  sum[key].executed += Math.abs(Number(item.total) || 0) / 100;
  return sum;
}, { earns: { planned: 0, executed: 0 }, expenses: { planned: 0, executed: 0 }, investments: { planned: 0, executed: 0 } }));
const chartData = computed(() => ({
  labels: ["Receitas", "Despesas", "Investimentos"], datasets: [
    { label: "Planejado", data: [totals.value.earns.planned, totals.value.expenses.planned, totals.value.investments.planned], backgroundColor: "#a5b4fc", borderRadius: 6 },
    { label: "Realizado", data: [totals.value.earns.executed, totals.value.expenses.executed, totals.value.investments.executed], backgroundColor: "#4f46e5", borderRadius: 6 },
  ]
}));
const chartOptions = { responsive: true, maintainAspectRatio: false, plugins: { legend: { position: "bottom", labels: { usePointStyle: true } }, tooltip: { callbacks: { label: (context) => ` ${context.dataset.label}: ${currencyBRL(context.raw)}` } } }, scales: { x: { grid: { display: false } }, y: { beginAtZero: true, ticks: { callback: (value) => currencyBRL(value) } } } };
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
