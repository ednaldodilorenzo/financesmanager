<template>
  <section class="transactions-report-page py-3">
    <div class="d-flex flex-wrap align-items-center justify-content-between gap-2 mb-3">
      <div>
        <h1 class="h4 fw-bold mb-1">Relatórios</h1>
        <p class="text-body-secondary small mb-0">Filtre seus lançamentos por período, descrição ou tags</p>
      </div>
      <button class="btn btn-outline-success" type="button" :disabled="!transactions.length" @click="exportToCSV"
        title="Exportar CSV">
        <i class="bi bi-download me-1"></i>CSV
      </button>
    </div>

    <div class="card border-0 shadow-sm mb-3">
      <form class="card-body p-3" @submit.prevent="search">
        <div class="d-flex flex-wrap align-items-end gap-3">
          <div>
            <label class="form-label small fw-semibold mb-1 d-block">Período</label>
            <div class="btn-group" role="group" aria-label="Tipo de período">
              <input id="periodMonth" v-model="periodMode" class="btn-check" type="radio" value="month" />
              <label class="btn btn-outline-primary" for="periodMonth">Mês/Ano</label>
              <input id="periodRange" v-model="periodMode" class="btn-check" type="radio" value="range" />
              <label class="btn btn-outline-primary" for="periodRange">Intervalo</label>
            </div>
          </div>

          <Calendar v-if="periodMode === 'month'" v-model="currentDate" />
          <template v-else>
            <div class="date-filter">
              <label class="form-label small fw-semibold mb-1" for="iptStartDate">Data inicial</label>
              <input id="iptStartDate" v-model="startDate" class="form-control" type="date" :max="endDate || undefined" />
            </div>
            <div class="date-filter">
              <label class="form-label small fw-semibold mb-1" for="iptEndDate">Data final</label>
              <input id="iptEndDate" v-model="endDate" class="form-control" type="date" :min="startDate || undefined" />
            </div>
          </template>

          <div class="description-filter">
            <label class="form-label small fw-semibold mb-1" for="iptDescription">Descrição</label>
            <input id="iptDescription" v-model="description" class="form-control" type="text"
              placeholder="Buscar pela descrição" />
          </div>

          <div class="tag-filter">
            <bootstrap-select-tag 
                v-model="tags" 
                :options="searchTags" 
                label="Tags" 
                id="iptTags"
                name="value" />
          </div>

          <button class="btn btn-primary px-3" type="submit">
            <i class="bi bi-search me-1"></i>Filtrar
          </button>
        </div>
      </form>
    </div>

    <div class="row g-4 align-items-start">
      <div class="col-xl-8">
        <div class="card border-0 shadow-sm transaction-list-card">
          <div class="card-header bg-white border-bottom px-3 py-3 d-flex align-items-center justify-content-between">
            <h2 class="h5 fw-bold mb-0 text-capitalize">{{ periodTitle }}</h2>
            <span class="badge rounded-pill text-bg-light">{{ transactions.length }} lançamentos</span>
          </div>

          <div v-if="groupedTransactions.length">
            <section v-for="group in groupedTransactions" :key="group.key" class="transaction-group">
              <div class="day-header d-flex justify-content-between align-items-center px-3 py-2">
                <strong>{{ group.label }}</strong>
                <strong :class="amountClass(group.total)">{{ currencyBRL(group.total) }}</strong>
              </div>

              <article v-for="item in group.items" :key="item.id"
                class="transaction-row d-flex align-items-center gap-3 px-3 py-3">
                <div :class="['transaction-icon', iconTone(item.categoryType)]">
                  <i :class="transactionIcon(item.categoryType)"></i>
                </div>
                <div class="min-w-0 flex-grow-1">
                  <div class="fw-semibold text-truncate">{{ item.description }}</div>
                  <div class="small text-body-secondary text-truncate">
                    {{ item.category }} <span class="mx-1">•</span> {{ item.account }}
                  </div>
                </div>
                <strong class="transaction-value text-nowrap" :class="amountClass(item.value)">
                  {{ currencyBRL(item.value) }}
                </strong>
              </article>
            </section>

            <div class="report-footer d-flex justify-content-between align-items-center px-3 py-3">
              <strong>Saldo final</strong>
              <strong class="h5 mb-0" :class="amountClass(finalBalance)">{{ currencyBRL(finalBalance) }}</strong>
            </div>
          </div>
          <div v-else class="text-center py-5 px-3">
            <i class="bi bi-receipt display-5 text-body-tertiary"></i>
            <h3 class="h6 mt-3 mb-1">Nenhum lançamento encontrado</h3>
            <p class="small text-body-secondary mb-0">Tente alterar os filtros selecionados.</p>
          </div>
        </div>
      </div>

      <div class="col-xl-4">
        <SummaryData :data="summaryCards" layout="sidebar" />
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed, ref } from "vue";
import { useToast } from "vue-toastification";
import Calendar from "@/components/bootstrap-calendar.vue";
import BootstrapSelectTag from "@/components/bootstrap-select-tag.vue";
import SummaryData from "@/components/summary-data.vue";
import transactionService from "./transaction.service";
import tagService from "@/views/tag/tag.service";
import { useLoadingScreen } from "@/components/loading/useLoadingScreen";
import { currencyBRL } from "@/components/filters/currency.filter";
import { formatDateUTC } from "@/utils/date";

const toast = useToast();
const loading = useLoadingScreen();

const transactions = ref([]);
const periodMode = ref("month");
const currentDate = ref(new Date());
const startDate = ref("");
const endDate = ref("");
const description = ref("");
const tags = ref([]);
// Title of the period that produced the current result, not of the pending filter.
const periodTitle = ref("");

const summary = computed(() => transactions.value.reduce((acc, item) => {
  if (item.categoryType === "R") acc.earns += item.value;
  if (item.categoryType === "D") acc.expenses += item.value;
  if (item.categoryType === "I") acc.investments -= item.value;
  return acc;
}, { earns: 0, expenses: 0, investments: 0 }));

const finalBalance = computed(() => summary.value.earns + summary.value.expenses - summary.value.investments);

const summaryCards = computed(() => [
  { title: "Saldo final", value: finalBalance.value, icon: "bi-wallet2" },
  { title: "Saldo sem investimentos", value: summary.value.earns + summary.value.expenses, icon: "bi-wallet2" },
  { title: "Receitas", value: summary.value.earns, icon: "bi-arrow-down-left", tone: "success" },
  { title: "Despesas", value: Math.abs(summary.value.expenses), icon: "bi-arrow-up-right", tone: "danger" },
  { title: "Investimentos", value: summary.value.investments, icon: "bi-graph-up-arrow", tone: "primary" },
]);

const dateKey = (value) => {
  if (value instanceof Date) return value.toISOString().slice(0, 10);
  return String(value).slice(0, 10);
};

const groupedTransactions = computed(() => {
  const groups = new Map();
  [...transactions.value]
    .sort((a, b) => new Date(b.paymentDate) - new Date(a.paymentDate))
    .forEach((item) => {
      const key = dateKey(item.paymentDate);
      const date = new Date(`${key}T12:00:00`);
      if (!groups.has(key)) groups.set(key, {
        key,
        label: new Intl.DateTimeFormat("pt-BR", { day: "2-digit", month: "short", year: "numeric", weekday: "long" }).format(date),
        total: 0,
        items: [],
      });
      const group = groups.get(key);
      group.items.push(item);
      group.total += item.value;
    });
  return [...groups.values()];
});

const amountClass = (value) => value > 0 ? "text-primary" : value < 0 ? "text-danger" : "text-body";
const iconTone = (type) => ({ R: "is-income", D: "is-expense", I: "is-investment" }[type] || "");
const transactionIcon = (type) => ({ R: "bi bi-arrow-down-left", D: "bi bi-bag", I: "bi bi-graph-up-arrow" }[type] || "bi bi-receipt");

const searchTags = (filter) => tagService
  .findAll({ filter })
  .then((resp) => resp.map((item) => item.tag));

function mapTransactions(list) {
  return list.map((item) => ({
    ...item,
    description: item.account.type === "C" ? `${item.description} - (${formatDateUTC(item.transactionDate, "dd/MM/yyyy")})` : item.description,
    category: item.category.name,
    categoryId: item.category.id,
    categoryType: item.category.type,
    account: item.account.name,
  }));
}

function buildParams() {
  // requester concatenates params as-is, so values are encoded here.
  const params = {};
  if (periodMode.value === "month") {
    params.month = currentDate.value.getMonth() + 1;
    params.year = currentDate.value.getFullYear();
  } else {
    params.startDate = encodeURIComponent(`${startDate.value}T00:00:00Z`);
    params.endDate = encodeURIComponent(`${endDate.value}T23:59:59Z`);
  }
  if (description.value.trim()) params.description = encodeURIComponent(description.value.trim());
  if (tags.value.length) params.tags = tags.value.map(encodeURIComponent).join(",");
  return params;
}

function buildPeriodTitle() {
  if (periodMode.value === "month") {
    return new Intl.DateTimeFormat("pt-BR", { month: "long", year: "numeric" }).format(currentDate.value);
  }
  return `${formatDateUTC(`${startDate.value}T00:00:00Z`, "dd/MM/yyyy")} a ${formatDateUTC(`${endDate.value}T00:00:00Z`, "dd/MM/yyyy")}`;
}

async function search() {
  if (periodMode.value === "range") {
    if (!startDate.value || !endDate.value) {
      toast.warning("Informe a data inicial e a data final.", { position: "top-center" });
      return;
    }
    if (startDate.value > endDate.value) {
      toast.warning("A data inicial deve ser anterior à data final.", { position: "top-center" });
      return;
    }
  }

  loading.show();
  try {
    const response = await transactionService.findAll(buildParams());
    transactions.value = mapTransactions(response);
    periodTitle.value = buildPeriodTitle();
  } finally {
    loading.hide();
  }
}

function exportToCSV() {
  const escape = (value) => `"${String(value ?? "").replaceAll('"', '""')}"`;
  const rows = transactions.value.map((item) => [formatDateUTC(item.paymentDate, "dd/MM/yyyy"), item.description, item.value, item.category, item.account].map(escape).join(";"));
  const footer = [escape("Saldo final"), "", escape(finalBalance.value), "", ""].join(";");
  const blob = new Blob(["﻿Data;Descrição;Valor;Categoria;Conta\n", ...rows.map((row) => `${row}\n`), `${footer}\n`], { type: "text/csv;charset=utf-8" });
  const link = Object.assign(document.createElement("a"), { href: URL.createObjectURL(blob), download: "relatorio-transacoes.csv" });
  link.click();
  URL.revokeObjectURL(link.href);
}

search();
</script>

<style scoped>
.transactions-report-page {
  --surface-muted: #f3f4f6;
}

.transaction-list-card {
  border-radius: .85rem;
}

.transaction-list-card .card-header {
  border-radius: .85rem .85rem 0 0;
}

.date-filter {
  width: min(100%, 170px);
}

.description-filter {
  width: min(100%, 240px);
}

.tag-filter {
  width: min(100%, 280px);
}

.tag-filter :deep(.container) {
  padding: 0;
}

.tag-filter :deep(.form-label) {
  font-size: .875em;
  font-weight: 600;
  margin-bottom: .25rem;
}

.day-header {
  background: var(--surface-muted);
  color: var(--bs-secondary-color);
  font-size: .875rem;
}

.report-footer {
  background: var(--surface-muted);
  border-radius: 0 0 .85rem .85rem;
}

.transaction-row {
  border-bottom: 1px solid var(--bs-border-color-translucent);
  transition: background-color .15s ease;
}

.transaction-row:hover {
  background: rgba(var(--bs-primary-rgb), .035);
}

.transaction-icon {
  width: 2.25rem;
  height: 2.25rem;
  display: grid;
  place-items: center;
  flex: 0 0 auto;
  border-radius: .65rem;
}

.transaction-icon:not(.is-income):not(.is-expense):not(.is-investment) {
  background: var(--bs-light);
  color: var(--bs-secondary);
}

.transaction-value {
  min-width: 7.5rem;
  text-align: right;
}

@media (max-width: 575.98px) {

  .date-filter,
  .description-filter,
  .tag-filter {
    width: 100%;
  }

  .transaction-row {
    gap: .65rem !important;
  }

  .transaction-value {
    min-width: auto;
    font-size: .9rem;
  }
}
</style>
