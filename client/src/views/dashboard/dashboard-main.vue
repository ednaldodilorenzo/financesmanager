<template>
  <section class="dashboard-page py-3">
    <header class="d-flex flex-wrap align-items-center justify-content-between gap-3 mb-4">
      <div>
        <h1 class="h4 fw-bold mb-1">Dashboard</h1>
        <p class="small text-body-secondary mb-0">Visão geral da sua vida financeira</p>
      </div>
      <div class="period-card d-flex align-items-center gap-3 bg-white shadow-sm px-3 py-2">
        <div class="d-none d-sm-block"><span class="small text-body-secondary d-block">Período</span><strong
            class="text-capitalize">{{ periodLabel }}</strong></div>
        <Calendar @date-change="onDateChange" />
      </div>
    </header>

    <div class="row row-cols-1 row-cols-sm-2 row-cols-xl-3 row-cols-xxl-5 g-3 mb-4">
      <div v-for="metric in metrics" :key="metric.title" class="col">
        <div class="metric-card card border-0 shadow-sm h-100">
          <div class="card-body p-3">
            <div class="d-flex align-items-start justify-content-between gap-2">
              <div class="min-w-0">
                <p class="small fw-semibold text-body-secondary mb-2">{{ metric.title }}</p>
                <h2 :class="['h5 fw-bold mb-0 text-nowrap', metric.valueClass]">{{ currencyBRL(metric.value) }}</h2>
              </div>
              <span :class="['metric-icon', metric.tone]"><i :class="['bi', metric.icon]"></i></span>
            </div>
            <p class="small text-body-secondary mt-3 mb-0">{{ metric.description }}</p>
          </div>
        </div>
      </div>
    </div>

    <div class="row g-4">
      <div class="col-xl-5">
        <div class="card chart-card border-0 shadow-sm h-100">
          <div class="card-header bg-white border-bottom p-3">
            <h2 class="h6 fw-bold mb-1">Distribuição do mês</h2><span class="small text-body-secondary">Receitas,
              despesas e investimentos</span>
          </div>
          <div class="card-body chart-body">
            <div class="pie-chart-container">
              <PieChart v-if="hasTypeData" :chart-data="chartData" :options="doughnutOptions" />
              <ChartEmpty v-else />
            </div>
          </div>
        </div>
      </div>
      <div class="col-xl-7">
        <DayByDayGraph :transactions-list="transactionsList" :date="currentDate" />
      </div>
      <div class="col-xl-6">
        <div class="card chart-card border-0 shadow-sm h-100">
          <div class="card-header bg-white border-bottom p-3">
            <h2 class="h6 fw-bold mb-1">Despesas por categoria</h2><span class="small text-body-secondary">Onde o
              dinheiro foi utilizado</span>
          </div>
          <div class="card-body chart-body">
            <div class="pie-chart-container">
              <PieChart v-if="chartExpensesData.labels.length" :chart-data="chartExpensesData"
                :options="doughnutOptions" />
              <ChartEmpty v-else />
            </div>
          </div>
        </div>
      </div>
      <div class="col-xl-6">
        <div class="card chart-card border-0 shadow-sm h-100">
          <div class="card-header bg-white border-bottom p-3">
            <h2 class="h6 fw-bold mb-1">Receitas por categoria</h2><span class="small text-body-secondary">Origem das
              entradas do mês</span>
          </div>
          <div class="card-body chart-body">
            <div class="pie-chart-container">
              <PieChart v-if="chartEarnsData.labels.length" :chart-data="chartEarnsData" :options="doughnutOptions" />
              <ChartEmpty v-else />
            </div>
          </div>
        </div>
      </div>
      <div class="col-xl-6">
        <EarnsExpensesYear :date="currentDate" :data-list="transactionYearSummaryList" />
      </div>
      <div class="col-xl-6">
        <PlannedExecutedType :data-list="plannedList" />
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed, defineComponent, h, ref } from "vue";
import { PieChart } from "vue-chart-3";
import { Chart, registerables } from "chart.js";
import { useRouter } from "vue-router";
import Calendar from "@/components/bootstrap-calendar.vue";
import transactionService from "@/views/transaction/transaction.service";
import planningService from "../planning/planning.service";
import { useLoadingScreen } from "@/components/loading/useLoadingScreen";
import { currencyBRL } from "@/components/filters/currency.filter";
import DayByDayGraph from "./daybyday-graph.vue";
import PlannedExecutedType from "./planned-executed-type.vue";
import EarnsExpensesYear from "./earns-expenses-year.vue";

Chart.register(...registerables);
const ChartEmpty = defineComponent({ setup: () => () => h("div", { class: "h-100 d-flex flex-column align-items-center justify-content-center text-body-secondary" }, [h("i", { class: "bi bi-bar-chart fs-2" }), h("span", { class: "small mt-2" }, "Sem dados para este período")]) });
const loading = useLoadingScreen();
const router = useRouter();
const currentDate = ref(new Date());
const transactionsList = ref([]);
const plannedList = ref([]);
const transactionYearSummaryList = ref([]);
const palette = ["#4f46e5", "#06b6d4", "#f59e0b", "#8b5cf6", "#10b981", "#ef4444", "#3b82f6", "#ec4899"];

const periodLabel = computed(() => new Intl.DateTimeFormat("pt-BR", { month: "long", year: "numeric" }).format(currentDate.value));
const totals = computed(() => transactionsList.value.reduce((sum, item) => {
  const value = Math.abs(Number(item.value) || 0) / 100;
  if (item.category?.type === "R") sum.earns += value;
  if (item.category?.type === "D") sum.expenses += value;
  if (item.category?.type === "I") sum.investments += value;
  return sum;
}, { earns: 0, expenses: 0, investments: 0 }));
const monthBalance = computed(() => totals.value.earns - totals.value.expenses);
const availableBalance = computed(() => monthBalance.value - totals.value.investments);
const metrics = computed(() => [
  { title: "Receitas", value: totals.value.earns, description: "Entradas no mês", icon: "bi-arrow-down-left", tone: "is-income", valueClass: "text-success" },
  { title: "Despesas", value: totals.value.expenses, description: "Saídas no mês", icon: "bi-arrow-up-right", tone: "is-expense", valueClass: "text-danger" },
  { title: "Saldo do mês", value: monthBalance.value, description: "Receitas menos despesas", icon: "bi-wallet2", tone: monthBalance.value < 0 ? "is-expense" : "is-balance", valueClass: monthBalance.value < 0 ? "text-danger" : "text-primary" },
  { title: "Valor investido", value: totals.value.investments, description: "Aportes realizados", icon: "bi-graph-up-arrow", tone: "is-investment", valueClass: "text-primary" },
  { title: "Saldo disponível", value: availableBalance.value, description: "Saldo após investimentos", icon: "bi-piggy-bank", tone: availableBalance.value < 0 ? "is-expense" : "is-available", valueClass: availableBalance.value < 0 ? "text-danger" : "text-success" },
]);
const valuesByCategory = computed(() => transactionsList.value.reduce((result, item) => {
  const type = item.category?.type;
  if (!["R", "D"].includes(type)) return result;
  const name = item.category?.name || "Sem categoria";
  const existing = result.find((entry) => entry.category === name && entry.type === type);
  const value = Math.abs(Number(item.value) || 0) / 100;
  if (existing) existing.value += value; else result.push({ category: name, value, type });
  return result;
}, []));
const chartData = computed(() => ({ labels: ["Receitas", "Despesas", "Investimentos"], datasets: [{ data: [totals.value.earns, totals.value.expenses, totals.value.investments], backgroundColor: ["#10b981", "#ef4444", "#4f46e5"], borderWidth: 0 }] }));
const hasTypeData = computed(() => chartData.value.datasets[0].data.some((value) => value > 0));
const categoryChart = (type, label) => computed(() => {
  const data = valuesByCategory.value.filter((item) => item.type === type);
  return { labels: data.map((item) => item.category), datasets: [{ label, data: data.map((item) => item.value), backgroundColor: data.map((_, index) => palette[index % palette.length]), borderWidth: 0 }] };
});
const chartExpensesData = categoryChart("D", "Despesas");
const chartEarnsData = categoryChart("R", "Receitas");
const doughnutOptions = { responsive: true, maintainAspectRatio: false, cutout: "62%", plugins: { legend: { position: "bottom", labels: { usePointStyle: true, padding: 18 } }, tooltip: { callbacks: { label: (context) => ` ${context.label}: ${currencyBRL(context.raw)}` } } } };

async function loadData() {
  loading.show();
  const params = { month: currentDate.value.getMonth() + 1, year: currentDate.value.getFullYear() };
  try {
    const [transactions, planning, yearly] = await Promise.all([transactionService.findAll(params), planningService.findAll(params), transactionService.getAllByYear(params.year)]);
    transactionsList.value = transactions.data;
    plannedList.value = planning.data;
    transactionYearSummaryList.value = yearly.data;
  } catch (error) { router.push({ name: "denied" }); }
  finally { loading.hide(); }
}
function onDateChange(date) { currentDate.value = date; loadData(); }
loadData();
</script>

<style scoped>
.period-card,
.metric-card,
.chart-card {
  border-radius: .85rem
}

.metric-card {
  transition: transform .15s ease
}

.metric-card:hover {
  transform: translateY(-2px)
}

.metric-icon {
  width: 2.5rem;
  height: 2.5rem;
  display: grid;
  place-items: center;
  flex: 0 0 auto;
  border-radius: .75rem
}

.is-income {
  color: var(--bs-success);
  background: rgba(var(--bs-success-rgb), .1)
}

.is-expense {
  color: var(--bs-danger);
  background: rgba(var(--bs-danger-rgb), .1)
}

.is-investment,
.is-balance {
  color: var(--bs-primary);
  background: rgba(var(--bs-primary-rgb), .1)
}

.is-available {
  color: var(--bs-success);
  background: rgba(var(--bs-success-rgb), .1)
}

.chart-card .card-header {
  border-radius: .85rem .85rem 0 0
}

.chart-body {
  height: 20rem;
  position: relative
}

.min-w-0 {
  min-width: 0
}

.pie-chart-container {
  position: relative;
  width: 100%;
  height: 260px;
  max-height: 260px;
}

.pie-chart-container :deep(canvas) {
  width: 100% !important;
  height: 100% !important;
  max-width: 100% !important;
  max-height: 260px !important;
}
</style>