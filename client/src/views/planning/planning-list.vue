<template>
  <section class="planning-page py-3">
    <header class="mb-3">
      <h1 class="h4 fw-bold mb-1">Planejamento</h1>
      <p class="small text-body-secondary mb-0">Compare o que foi planejado com o realizado no período</p>
    </header>

    <div class="card border-0 shadow-sm mb-4">
      <div class="card-body p-3">
        <div class="d-flex flex-wrap align-items-center justify-content-between gap-3">
          <div>
            <span class="small fw-semibold text-body-secondary d-block mb-1">Período analisado</span>
            <strong class="h5 text-capitalize mb-0">{{ periodLabel }}</strong>
          </div>
          <div class="d-flex flex-wrap align-items-center gap-3">
            <Calendar @date-change="onChangeDebounced" />
            <div class="period-toggle btn-group" role="group" aria-label="Tipo de período">
              <input id="period-month" v-model="type" type="radio" class="btn-check" value="M" @change="updateState" />
              <label class="btn btn-outline-primary" for="period-month">Mensal</label>
              <input id="period-year" v-model="type" type="radio" class="btn-check" value="Y" @change="updateState" />
              <label class="btn btn-outline-primary" for="period-year">Anual</label>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="row g-3 mb-4">
      <div v-for="card in overviewCards" :key="card.title" class="col-md-4">
        <div class="card border-0 shadow-sm h-100 overview-card">
          <div class="card-body p-3">
            <div class="d-flex justify-content-between align-items-start gap-3">
              <div>
                <p class="small fw-semibold text-body-secondary mb-2">{{ card.title }}</p>
                <h2 :class="['h5 fw-bold mb-1', `text-${card.tone}`]">{{ currencyBRL(card.executed) }}</h2>
                <span class="small text-body-secondary">de {{ currencyBRL(card.planned) }}</span>
              </div>
              <span :class="['overview-icon', `text-bg-${card.tone}`]"><i :class="['bi', card.icon]"></i></span>
            </div>
            <BootstrapPlanExecBar class="mt-3" :planned="card.planned" :executed="card.executed"
              :percent-divider="percentDivider" :tone="card.tone" compact />
          </div>
        </div>
      </div>
    </div>

    <div class="d-grid gap-4">
      <article v-for="section in sections" :key="section.type" class="card border-0 shadow-sm planning-card">
        <div class="card-header bg-white border-bottom p-3">
          <div class="d-flex flex-wrap align-items-center justify-content-between gap-3">
            <div class="d-flex align-items-center gap-3">
              <span :class="['section-icon', `text-bg-${section.tone}`]"><i :class="['bi', section.icon]"></i></span>
              <div>
                <h2 class="h5 fw-bold mb-1">{{ section.title }}</h2>
                <span class="small text-body-secondary">{{ section.items.length }} categorias</span>
              </div>
            </div>
            <div class="text-end">
              <span class="small text-body-secondary d-block">Total planejado</span>
              <strong :class="`text-${section.tone}`">{{ currencyBRL(section.summary.planned) }}</strong>
            </div>
          </div>
        </div>

        <div v-if="section.items.length">
          <div class="planning-total px-3 py-3 border-bottom">
            <BootstrapPlanExecBar :planned="section.summary.planned" :executed="section.summary.executed"
              :percent-divider="percentDivider" :tone="section.tone" />
          </div>
          <div v-for="item in section.items" :key="item.id || item.name" class="planning-row px-3 py-3">
            <div class="row align-items-center g-2 g-lg-3">
              <div class="col-lg-3">
                <div class="fw-semibold text-truncate">{{ item.name }}</div>
                <div class="small text-body-secondary d-lg-none">Planejado: {{ currencyBRL(item.planned) }}</div>
              </div>
              <div class="col-lg-7">
                <BootstrapPlanExecBar :planned="item.planned" :executed="item.executed"
                  :percent-divider="percentDivider" :tone="section.tone" />
              </div>
              <div class="col-lg-2 text-lg-end d-none d-lg-block">
                <strong :class="`text-${section.tone}`">{{ currencyBRL(item.planned) }}</strong>
              </div>
            </div>
          </div>
        </div>
        <div v-else class="text-center py-5 px-3">
          <i class="bi bi-inbox fs-2 text-body-tertiary"></i>
          <p class="small text-body-secondary mt-2 mb-0">Nenhuma categoria neste período.</p>
        </div>
      </article>
    </div>
  </section>
</template>

<script setup>
import { computed, ref } from "vue";
import Calendar from "@/components/bootstrap-calendar.vue";
import BootstrapPlanExecBar from "@/components/bootstrap-planexec-bar.vue";
import planningService from "./planning.service";
import { debounce } from "@/utils/support";
import { useLoadingScreen } from "@/components/loading/useLoadingScreen";
import { currencyBRL } from "@/components/filters/currency.filter";

const loading = useLoadingScreen();
const expensesList = ref([]);
const earnsList = ref([]);
const investmentsList = ref([]);
const fullList = ref([]);
const type = ref("M");
const currentDate = ref(new Date());

const summarize = (list) => list.reduce((total, item) => ({
  planned: total.planned + (Number(item.planned) || 0),
  executed: total.executed + (Number(item.executed) || 0),
}), { planned: 0, executed: 0 });

const investmentsSummary = computed(() => summarize(investmentsList.value));
const earnsSummary = computed(() => summarize(earnsList.value));
const expensesSummary = computed(() => summarize(expensesList.value));

const percentDivider = computed(() => type.value === "Y" ? (currentDate.value.getMonth() + 1) / 12 : 0);
const periodLabel = computed(() => new Intl.DateTimeFormat("pt-BR", type.value === "M"
  ? { month: "long", year: "numeric" }
  : { year: "numeric" }
).format(currentDate.value));

const sections = computed(() => [
  { type: "I", title: "Investimentos", items: investmentsList.value, summary: investmentsSummary.value, tone: "primary", icon: "bi-graph-up-arrow" },
  { type: "R", title: "Receitas", items: earnsList.value, summary: earnsSummary.value, tone: "success", icon: "bi-arrow-down-left" },
  { type: "D", title: "Despesas", items: expensesList.value, summary: expensesSummary.value, tone: "danger", icon: "bi-arrow-up-right" },
]);

const overviewCards = computed(() => sections.value.map((section) => ({
  title: section.title,
  planned: section.summary.planned,
  executed: section.summary.executed,
  tone: section.tone,
  icon: section.icon,
})));

function updateState() {
  const monthly = type.value === "M";
  const mapped = fullList.value.map((item) => {
    const planned = Math.abs(Number(item.planned) || 0) / (monthly ? 12 : 1);
    const executed = Math.abs(Number(monthly ? item.total : item.accumulated) || 0);
    return { ...item, planned, executed };
  });
  expensesList.value = mapped.filter(({ type }) => type === "D");
  earnsList.value = mapped.filter(({ type }) => type === "R");
  investmentsList.value = mapped.filter(({ type }) => type === "I");
}

async function getData(month, year) {
  loading.show();
  try {
    const response = await planningService.findAll({ month, year });
    fullList.value = response.data;
    updateState();
  } finally {
    loading.hide();
  }
}

const onChangeDebounced = debounce((date) => {
  currentDate.value = date;
  getData(date.getMonth() + 1, date.getFullYear());
}, 500);

getData(currentDate.value.getMonth() + 1, currentDate.value.getFullYear());
</script>

<style scoped>
.period-toggle .btn {
  min-width: 5.5rem;
}

.overview-card,
.planning-card {
  border-radius: .85rem;
}

.overview-icon,
.section-icon {
  display: grid;
  place-items: center;
  flex: 0 0 auto;
  --bs-bg-opacity: .12;
}

.overview-icon {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: .75rem;
}

.section-icon {
  width: 2.75rem;
  height: 2.75rem;
  border-radius: .8rem;
}

.planning-card .card-header {
  border-radius: .85rem .85rem 0 0;
}

.planning-total {
  background: var(--bs-tertiary-bg);
}

.planning-row {
  border-bottom: 1px solid var(--bs-border-color-translucent);
  transition: background-color .15s ease;
}

.planning-row:last-child {
  border-bottom: 0;
}

.planning-row:hover {
  background: rgba(var(--bs-primary-rgb), .025);
}

@media (max-width: 575.98px) {
  .period-toggle {
    width: 100%;
  }

  .period-toggle .btn {
    flex: 1;
  }
}
</style>