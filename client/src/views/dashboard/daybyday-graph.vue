<template>
  <div class="card chart-card border-0 shadow-sm h-100">
    <div class="card-header bg-white border-bottom p-3">
      <h2 class="h6 fw-bold mb-1">Evolução diária</h2><span class="small text-body-secondary">Receitas e despesas
        acumuladas no mês</span>
    </div>
    <div class="card-body chart-body">
      <div class="chart-container">
        <LineChart v-if="transactionsList.length" :chart-data="chartData" :options="chartOptions" />
        <div v-else class="h-100 d-flex flex-column align-items-center justify-content-center text-body-secondary"><i
            class="bi bi-graph-up fs-2"></i><span class="small mt-2">Sem dados para este período</span></div>
      </div>
    </div>
  </div>
</template>
<script setup>
import { computed } from "vue";
import { LineChart } from "vue-chart-3";
import { currencyBRL } from "@/components/filters/currency.filter";

const props = defineProps({ transactionsList: { type: Array, default: () => [] }, date: { type: Date, default: () => new Date() } });
const days = computed(() => {
  const total = new Date(props.date.getFullYear(), props.date.getMonth() + 1, 0).getDate();
  const today = new Date();
  const isCurrent = props.date.getFullYear() === today.getFullYear() && props.date.getMonth() === today.getMonth();
  const limit = isCurrent ? today.getDate() : total;
  return Array.from({ length: limit }, (_, index) => index + 1);
});
const evolution = computed(() => {
  const daily = new Map(days.value.map((day) => [day, { earns: 0, expenses: 0 }]));
  props.transactionsList.forEach((item) => {
    const day = new Date(item.paymentDate).getUTCDate();
    const entry = daily.get(day);
    if (!entry) return;
    const value = Math.abs(Number(item.value) || 0) / 100;
    if (item.category?.type === "R") entry.earns += value;
    if (item.category?.type === "D") entry.expenses += value;
  });
  let earns = 0, expenses = 0;
  return days.value.map((day) => { const item = daily.get(day); earns += item.earns; expenses += item.expenses; return { day, earns, expenses }; });
});
const chartData = computed(() => ({
  labels: days.value, datasets: [
    { label: "Receitas", data: evolution.value.map((item) => item.earns), borderColor: "#10b981", backgroundColor: "rgba(16,185,129,.1)", fill: true, tension: .35, pointRadius: 1.5 },
    { label: "Despesas", data: evolution.value.map((item) => item.expenses), borderColor: "#ef4444", backgroundColor: "rgba(239,68,68,.08)", fill: true, tension: .35, pointRadius: 1.5 },
  ]
}));
const chartOptions = { responsive: true, maintainAspectRatio: false, interaction: { mode: "index", intersect: false }, plugins: { legend: { position: "bottom", labels: { usePointStyle: true } }, tooltip: { callbacks: { label: (context) => ` ${context.dataset.label}: ${currencyBRL(context.raw)}` } } }, scales: { x: { grid: { display: false } }, y: { beginAtZero: true, ticks: { callback: (value) => currencyBRL(value) } } } };
</script>
