<template>
  <section class="transactions-page py-3">
    <div class="d-flex flex-wrap align-items-center justify-content-between gap-2 mb-3">
      <div>
        <h1 class="h4 fw-bold mb-1">Transações</h1>
        <p class="text-body-secondary small mb-0">Acompanhe seus lançamentos do período</p>
      </div>
      <button class="btn btn-primary px-3" type="button" @click="onNewClicked">
        <i class="bi bi-plus-lg me-1"></i>Adicionar lançamento
      </button>
    </div>

    <div class="card border-0 shadow-sm mb-3">
      <div class="card-body p-3">
        <div class="d-flex flex-wrap align-items-center gap-3">
          <Calendar class="me-auto" @date-change="onChangeDebounced" />
          <select v-model="selectedType" class="form-select filter-control" aria-label="Tipo">
            <option value="">Todos os tipos</option>
            <option value="R">Receitas</option>
            <option value="D">Despesas</option>
            <option value="I">Investimentos</option>
          </select>
          <select v-model="selectedAccount" class="form-select filter-control" aria-label="Conta">
            <option value="">Todas as contas</option>
            <option v-for="account in accounts" :key="account.id" :value="account.name">
              {{ account.name }}
            </option>
          </select>
          <div class="category-filter">
            <BootstrapSearcheableSelect displayField="name" valueField="id" v-model="selectedCategory"
              :options="categories" />
          </div>
          <button class="btn btn-outline-success" type="button" @click="exportToCSV" title="Exportar CSV">
            <i class="bi bi-download me-1"></i>CSV
          </button>
        </div>
      </div>
    </div>

    <div class="row g-4 align-items-start">
      <div class="col-xl-8">
        <div class="card border-0 shadow-sm transaction-list-card">
          <div class="card-header bg-white border-bottom px-3 py-3 d-flex align-items-center justify-content-between">
            <h2 class="h5 fw-bold mb-0 text-capitalize">{{ monthTitle }}</h2>
            <span class="badge rounded-pill text-bg-light">{{ filteredItems.length }} lançamentos</span>
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
                <div class="dropdown">
                  <button class="btn btn-sm btn-light rounded-circle" type="button" aria-label="Ações"
                    aria-haspopup="true" :aria-expanded="openDropdownId === item.id"
                    @click.stop="toggleDropdown(item.id)">
                    <i class="bi bi-three-dots-vertical"></i>
                  </button>
                  <ul class="dropdown-menu dropdown-menu-end shadow-sm border-0"
                    :class="{ show: openDropdownId === item.id }" @click.stop>
                    <li><button class="dropdown-item" type="button" @click="handleEdit(item)"><i
                          class="bi bi-pencil me-2"></i>Editar</button></li>
                    <li><button class="dropdown-item text-danger" type="button" @click="handleDelete(item)"><i
                          class="bi bi-trash me-2"></i>Excluir</button></li>
                  </ul>
                </div>
              </article>
            </section>
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
import { computed, onBeforeUnmount, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { useToast } from "vue-toastification";
import Calendar from "@/components/bootstrap-calendar.vue";
import BootstrapSearcheableSelect from "@/components/bootstrap-searcheable-select.vue";
import SummaryData from "@/components/summary-data.vue";
import TransactionChange from "./transaction-change.vue";
import transactionService from "./transaction.service";
import categoryService from "../category/category.service";
import accountService from "../account/account.service";
import { useModalScreen } from "@/components/modal/use-modal-screen";
import { useLoadingScreen } from "@/components/loading/useLoadingScreen";
import { useDialogScreen } from "@/components/dialog/use-dialog-screen";
import { currencyBRL } from "@/components/filters/currency.filter";
import { formatDateUTC } from "@/utils/date";
import { debounce } from "@/utils/support";

const router = useRouter();
const toast = useToast();
const modal = useModalScreen(TransactionChange);
const loading = useLoadingScreen();
const dialog = useDialogScreen("Deseja realmente excluir a transação?", "Exclusão");

const transactions = ref([]);
const categories = ref([]);
const accounts = ref([]);
const selectedType = ref("");
const selectedAccount = ref("");
const selectedCategory = ref(null);
const currentDate = ref(new Date());
const openDropdownId = ref(null);

const closeDropdown = () => { openDropdownId.value = null; };
const toggleDropdown = (id) => {
  openDropdownId.value = openDropdownId.value === id ? null : id;
};

onMounted(() => document.addEventListener("click", closeDropdown));
onBeforeUnmount(() => document.removeEventListener("click", closeDropdown));

const filteredItems = computed(() => transactions.value.filter((item) =>
  (!selectedType.value || item.categoryType === selectedType.value) &&
  (!selectedAccount.value || item.account === selectedAccount.value) &&
  (!selectedCategory.value || item.categoryId === selectedCategory.value.id)
));

const summary = computed(() => filteredItems.value.reduce((acc, item) => {
  if (item.categoryType === "R") acc.earns += item.value;
  if (item.categoryType === "D") acc.expenses += item.value;
  if (item.categoryType === "I") acc.investments += item.value;
  return acc;
}, { earns: 0, expenses: 0, investments: 0 }));

const summaryCards = computed(() => [
  { title: "Saldo do período", value: summary.value.earns + summary.value.expenses + summary.value.investments, icon: "bi-wallet2" },
  { title: "Receitas", value: summary.value.earns, icon: "bi-arrow-down-left", tone: "success" },
  { title: "Despesas", value: summary.value.expenses, icon: "bi-arrow-up-right", tone: "danger" },
  { title: "Investimentos", value: summary.value.investments, icon: "bi-graph-up-arrow", tone: "primary" },
]);

const monthTitle = computed(() => new Intl.DateTimeFormat("pt-BR", { month: "long", year: "numeric" }).format(currentDate.value));

const dateKey = (value) => {
  if (value instanceof Date) return value.toISOString().slice(0, 10);
  return String(value).slice(0, 10);
};

const groupedTransactions = computed(() => {
  const groups = new Map();
  [...filteredItems.value]
    .sort((a, b) => new Date(b.paymentDate) - new Date(a.paymentDate))
    .forEach((item) => {
      const key = dateKey(item.paymentDate);
      const date = new Date(`${key}T12:00:00`);
      if (!groups.has(key)) groups.set(key, {
        key,
        label: new Intl.DateTimeFormat("pt-BR", { day: "2-digit", weekday: "long" }).format(date),
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

async function loadData() {
  loading.show();
  try {
    const params = { month: currentDate.value.getMonth() + 1, year: currentDate.value.getFullYear() };
    const [tx, cats, accs] = await Promise.all([
      transactionService.findAll(params),
      categoryService.findAll({ paginate: false }),
      accountService.findAll(),
    ]);
    transactions.value = mapTransactions(tx.data);
    categories.value = cats.data;
    accounts.value = accs.data;
  } catch (error) {
    router.push({ name: "denied" });
  } finally {
    loading.hide();
  }
}

async function getList() {
  loading.show();
  try {
    const response = await transactionService.findAll({ month: currentDate.value.getMonth() + 1, year: currentDate.value.getFullYear() });
    transactions.value = mapTransactions(response.data);
  } finally { loading.hide(); }
}

const onChangeDebounced = debounce((date) => { currentDate.value = date; getList(); }, 500);

async function handleEdit(clicked) {
  closeDropdown();
  const saved = await modal.show({
    ...clicked,
    value: clicked.value.toLocaleString("pt-BR", { minimumFractionDigits: 2, maximumFractionDigits: 2 }),
    transactionDate: formatDateUTC(clicked.transactionDate, "yyyy-MM-dd"),
    paymentDate: formatDateUTC(clicked.paymentDate, "yyyy-MM-dd"),
  });
  if (saved) getList();
}

async function handleDelete(clicked) {
  closeDropdown();
  if (!await dialog.show()) return;
  loading.show();
  try {
    await transactionService.delete(clicked.id);
    transactions.value = transactions.value.filter(({ id }) => id !== clicked.id);
    toast.success("Transação excluída com sucesso!", { position: "top-center" });
  } finally { loading.hide(); }
}

async function onNewClicked() { if (await modal.show()) getList(); }

function exportToCSV() {
  const escape = (value) => `"${String(value ?? "").replaceAll('"', '""')}"`;
  const rows = filteredItems.value.map((item) => [formatDateUTC(item.paymentDate, "dd/MM/yyyy"), item.description, item.value, item.category, item.account].map(escape).join(";"));
  const blob = new Blob(["\uFEFFData;Descrição;Valor;Categoria;Conta\n", ...rows.map((row) => `${row}\n`)], { type: "text/csv;charset=utf-8" });
  const link = Object.assign(document.createElement("a"), { href: URL.createObjectURL(blob), download: "transacoes.csv" });
  link.click();
  URL.revokeObjectURL(link.href);
}

loadData();
</script>

<style scoped>
.transactions-page {
  --surface-muted: #f3f4f6;
}

.transaction-list-card {
  border-radius: .85rem;
}

.transaction-list-card .card-header {
  border-radius: .85rem .85rem 0 0;
}

.dropdown-menu.show {
  display: block;
  z-index: 1080;
}

.filter-control {
  width: min(100%, 180px);
}

.category-filter {
  width: min(100%, 220px);
}

.day-header {
  background: var(--surface-muted);
  color: var(--bs-secondary-color);
  font-size: .875rem;
}

.transaction-row {
  border-bottom: 1px solid var(--bs-border-color-translucent);
  transition: background-color .15s ease;
}

.transaction-row:hover {
  background: rgba(var(--bs-primary-rgb), .035);
}

.transaction-row:last-child {
  border-bottom: 0;
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

  .filter-control,
  .category-filter {
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