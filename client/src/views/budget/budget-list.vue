<template>
  <section class="budget-page py-3">
    <header class="d-flex flex-wrap align-items-center justify-content-between gap-2 mb-3">
      <div>
        <h1 class="h4 fw-bold mb-1">Orçamento</h1>
        <p class="small text-body-secondary mb-0">Planeje receitas, despesas e investimentos do ano</p>
      </div>
      <button class="btn btn-primary px-3" type="button" @click="clickNew">
        <i class="bi bi-plus-lg me-1"></i>Adicionar categoria
      </button>
    </header>

    <div class="card border-0 shadow-sm mb-3">
      <div class="card-body p-3">
        <div class="d-flex flex-wrap align-items-center justify-content-between gap-3">
          <div>
            <span class="small fw-semibold text-body-secondary d-block mb-1">Período do orçamento</span>
            <strong class="h5 mb-0">{{ currentYear }}</strong>
          </div>
          <Calendar @date-change="onChangeDebounced" :only-years="true" />
        </div>
      </div>
    </div>

    <SummaryData :data="summaryCards" />

    <div class="card border-0 shadow-sm mt-4 budget-table-card">
      <div
        class="card-header bg-white border-bottom p-3 d-flex flex-wrap align-items-center justify-content-between gap-2">
        <div>
          <h2 class="h5 fw-bold mb-1">Categorias planejadas</h2>
          <p class="small text-body-secondary mb-0">{{ filteredItems.length }} categorias no orçamento de {{ currentYear
            }}</p>
        </div>
        <button class="btn btn-sm btn-outline-primary" type="button" @click="clickNew">
          <i class="bi bi-plus-lg me-1"></i>Nova categoria
        </button>
      </div>

      <div v-if="filteredItems.length" class="table-responsive">
        <table class="table align-middle table-hover mb-0">
          <thead>
            <tr>
              <th>Categoria</th>
              <th class="text-end">Valor planejado</th>
              <th class="text-center actions-column">Ações</th>
            </tr>
          </thead>
          <tbody>
            <BudgetRecord v-for="(item, index) in filteredItems" :key="item.id || `new-${index}`" :budget-item="item"
              :options="allCategories" @cancel-item="cancelRecord" @delete-item="handleDelete"
              @change-item="handleChange" />
          </tbody>
        </table>
      </div>

      <div v-else class="text-center py-5 px-3">
        <span class="empty-icon"><i class="bi bi-pie-chart"></i></span>
        <h3 class="h6 mt-3 mb-1">Seu orçamento está vazio</h3>
        <p class="small text-body-secondary mb-3">Adicione uma categoria para começar o planejamento.</p>
        <button class="btn btn-primary btn-sm" type="button" @click="clickNew">
          <i class="bi bi-plus-lg me-1"></i>Adicionar categoria
        </button>
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed, ref } from "vue";
import Calendar from "@/components/bootstrap-calendar.vue";
import SummaryData from "@/components/summary-data.vue";
import { debounce } from "@/utils/support";
import { useLoadingScreen } from "@/components/loading/useLoadingScreen";
import BudgetRecord from "./budget-record.vue";
import BudgetItem from "./budget-item";
import budgetService from "./budget.service";
import categoryService from "@/views/category/category.service";

const loading = useLoadingScreen();
const currentDate = ref(new Date());
const filteredItems = ref([]);
const allCategories = ref([]);

const currentYear = computed(() => currentDate.value.getFullYear());

const summary = computed(() => filteredItems.value.reduce((total, item) => {
  const value = Number(item.value) || 0;
  if (item.category?.type === "R") total.earns += value;
  if (item.category?.type === "D") total.expenses += value;
  if (item.category?.type === "I") total.investments += value;
  return total;
}, { earns: 0, expenses: 0, investments: 0 }));

const ratio = (value) => summary.value.earns ? Math.abs(value / summary.value.earns) : 0;

const summaryCards = computed(() => {
  const balance = summary.value.earns - Math.abs(summary.value.expenses) - Math.abs(summary.value.investments);
  return [
    { title: "Total de receitas", value: summary.value.earns, icon: "bi-arrow-down-left", tone: "success" },
    { title: "Total de despesas", value: Math.abs(summary.value.expenses), percent: ratio(summary.value.expenses), percentMessage: "das receitas", icon: "bi-arrow-up-right", tone: "danger" },
    { title: "Saldo planejado", value: balance, percent: ratio(balance), percentMessage: "das receitas", icon: "bi-wallet2", tone: balance < 0 ? "danger" : "primary" },
    { title: "Investimento total", value: Math.abs(summary.value.investments), percent: ratio(summary.value.investments), percentMessage: "das receitas", icon: "bi-graph-up-arrow", tone: "primary" },
  ];
});

function handleDelete(item) {
  filteredItems.value = filteredItems.value.filter(({ id }) => id !== item.id);
}

function cancelRecord(item) {
  filteredItems.value = filteredItems.value.filter((value) => value !== item);
}

function handleChange(updatedItem, originalItem) {
  const index = filteredItems.value.findIndex((item) =>
    item === originalItem || (updatedItem.id && item.id === updatedItem.id)
  );
  if (index !== -1) {
    filteredItems.value.splice(index, 1, updatedItem);
  } else {
    filteredItems.value.unshift(updatedItem);
  }
}

async function getData() {
  loading.show();
  try {
    const [categoriesResponse, budgetResponse] = await Promise.all([
      categoryService.findAll({ paginate: false }),
      budgetService.findAll({ year: currentYear.value }),
    ]);
    allCategories.value = categoriesResponse.data;
    parseBudgetResponse(budgetResponse);
  } finally {
    loading.hide();
  }
}

function clickNew() {
  filteredItems.value.unshift(new BudgetItem({
    year: currentYear.value,
    value: 0,
    categoryId: 0,
  }, allCategories.value));
}

function parseBudgetResponse(response) {
  filteredItems.value = response.data.map((item) => new BudgetItem(item, allCategories.value));
}

async function getBudget(year) {
  loading.show();
  try {
    parseBudgetResponse(await budgetService.findAll({ year }));
  } finally {
    loading.hide();
  }
}

const onChangeDebounced = debounce((date) => {
  currentDate.value = date;
  getBudget(currentYear.value);
}, 500);

getData();
</script>

<style scoped>
.budget-table-card {
  border-radius: .85rem;
}

.budget-table-card .card-header {
  border-radius: .85rem .85rem 0 0;
}

thead th {
  padding: .8rem 1rem;
  background: var(--bs-tertiary-bg);
  border-bottom-width: 1px;
  color: var(--bs-secondary-color);
  font-size: .75rem;
  font-weight: 700;
  letter-spacing: .035em;
  text-transform: uppercase;
  white-space: nowrap;
}

.actions-column {
  width: 8.5rem;
}

.empty-icon {
  width: 3.25rem;
  height: 3.25rem;
  display: inline-grid;
  place-items: center;
  border-radius: 1rem;
  background: rgba(var(--bs-primary-rgb), .1);
  color: var(--bs-primary);
  font-size: 1.4rem;
}
</style>