<template>
  <div class="table-shell card border-0 shadow-sm overflow-hidden">
    <div v-if="showNav" class="card-header bg-white border-bottom p-3">
      <div class="d-flex flex-wrap align-items-center gap-2">
        <button v-if="showBack" class="btn btn-outline-secondary" type="button" @click="$emit('back-clicked')"
          aria-label="Voltar">
          <i class="bi bi-arrow-left"></i>
        </button>
        <button class="btn btn-primary" type="button" @click="$emit('new-clicked')">
          <i class="bi bi-plus-lg me-1"></i>Novo
        </button>
        <button v-if="showCSVButton" class="btn btn-outline-success" type="button" @click="csvHandler">
          <i class="bi bi-download me-1"></i>Exportar CSV
        </button>
        <div v-if="showFilter" class="search-box ms-lg-auto">
          <i class="bi bi-search"></i>
          <input v-model="searchQuery" class="form-control" type="search" placeholder="Pesquisar…"
            aria-label="Pesquisar" />
        </div>
      </div>
    </div>

    <div v-if="!filteredList.length" class="text-center py-5 px-3">
      <i class="bi bi-inbox display-5 text-body-tertiary"></i>
      <h2 class="h6 mt-3 mb-1">Nenhum dado encontrado</h2>
      <p class="small text-body-secondary mb-0">Altere a pesquisa ou adicione um novo registro.</p>
    </div>

    <div v-else class="table-responsive">
      <table :class="['table align-middle mb-0', { 'table-striped': striped, 'table-hover': hover }]">
        <thead>
          <tr>
            <th v-for="field in fields" :key="field.name" :class="fieldTitle(field).clazz">
              {{ fieldTitle(field).value }}
              <i v-if="sortColumns" class="bi bi-arrow-down-up ms-1 text-body-tertiary"></i>
            </th>
            <th v-if="actions.length" class="text-center actions-column">Ações</th>
          </tr>
        </thead>
        <tbody>
          <slot name="first-row" />
          <tr v-for="(item, index) in filteredList" :key="item.id ?? index" :class="item.clazz">
            <template v-for="field in fields" :key="field.name">
              <slot :name="`custom-td-${field.name}`" :item="item" :field="field">
                <td :style="cell(item, field).style" :class="cell(item, field).clazz">{{ cell(item, field).value }}</td>
              </slot>
            </template>
            <td v-if="actions.length" class="text-center text-nowrap">
              <button v-for="action in actions" :key="action.name" type="button"
                :class="['btn btn-sm btn-light rounded-circle ms-1', action.clazz]" :title="action.title"
                @click="action.handler(item)"><i :class="action.icon"></i></button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="showPagination && pageParams" class="card-footer bg-white border-top p-3">
      <nav aria-label="Paginação">
        <ul class="pagination pagination-sm justify-content-center mb-0">
          <li class="page-item" :class="{ disabled: pageParams.current === 1 }">
            <button class="page-link" @click="setPage(pageParams.current - 1)" aria-label="Anterior"><i
                class="bi bi-chevron-left"></i></button>
          </li>
          <li v-for="(page, index) in pageParams.items" :key="`${page}-${index}`" class="page-item"
            :class="{ active: page === pageParams.current, disabled: page === '…' }">
            <button class="page-link" @click="page !== '…' && setPage(page)">{{ page }}</button>
          </li>
          <li class="page-item" :class="{ disabled: pageParams.current === pageParams.max }">
            <button class="page-link" @click="setPage(pageParams.current + 1)" aria-label="Próxima"><i
                class="bi bi-chevron-right"></i></button>
          </li>
        </ul>
      </nav>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from "vue";
import { debounce } from "@/utils/support";

const props = defineProps({
  items: { type: [Array, Function], default: () => [] },
  csvHandler: { type: Function, default: () => { } },
  showCSVButton: { type: Boolean, default: false },
  pageSize: { type: Number, default: 10 },
  totalPages: { type: Number, default: 0 },
  currentPage: { type: Number, default: 1 },
  fields: { type: Array, default: () => [] },
  actions: { type: Array, default: () => [] },
  showFilter: { type: Boolean, default: true },
  showBack: { type: Boolean, default: false },
  showPagination: { type: Boolean, default: true },
  showNav: { type: Boolean, default: true },
  striped: { type: Boolean, default: false },
  hover: { type: Boolean, default: true },
  sortColumns: { type: Boolean, default: false },
});

const emit = defineEmits(["trigger-page", "new-clicked", "search-input", "back-clicked"]);
const searchQuery = ref("");
const filteredList = ref([]);
const totalRecords = ref(0);
const currentPage = ref(1);

const fieldTitle = (field) => typeof field.title === "object" ? field.title : { value: field.title, clazz: "" };
const cell = (item, field) => typeof item[field.name] === "object" && item[field.name] !== null
  ? item[field.name]
  : { value: item[field.name], clazz: "", style: "" };

async function loadData(page = 1, query = searchQuery.value) {
  const result = typeof props.items === "function"
    ? await props.items(page, props.pageSize, query)
    : getLocalItems(page, query);
  filteredList.value = result.items || [];
  totalRecords.value = result.total || 0;
  currentPage.value = result.page || page;
  emit("trigger-page", currentPage.value);
}

function getLocalItems(page, query) {
  const term = query.trim().toLocaleLowerCase("pt-BR");
  const matches = !term ? props.items : props.items.filter((item) =>
    Object.values(item).some((raw) => String(raw?.value ?? raw ?? "").toLocaleLowerCase("pt-BR").includes(term))
  );
  const start = (page - 1) * props.pageSize;
  return { items: props.showPagination ? matches.slice(start, start + props.pageSize) : matches, total: matches.length, page };
}

const setPage = (page) => {
  if (page < 1 || page > pageParams.value.max || page === currentPage.value) return;
  loadData(page);
};

const pageParams = computed(() => {
  const max = Math.ceil(totalRecords.value / props.pageSize);
  if (!max) return null;
  const current = currentPage.value;
  const visible = new Set([1, max, current - 1, current, current + 1]);
  const pages = [...visible].filter((page) => page >= 1 && page <= max).sort((a, b) => a - b);
  const items = [];
  pages.forEach((page, index) => {
    if (index && page - pages[index - 1] > 1) items.push("…");
    items.push(page);
  });
  return { current, max, items };
});

const search = debounce(() => { currentPage.value = 1; emit("search-input", searchQuery.value); loadData(1); }, 400);
watch(searchQuery, search);
watch(() => props.items, () => loadData(1), { deep: true });
onMounted(() => loadData());
</script>

<style scoped>
.table-shell {
  border-radius: .85rem;
}

.actions-column {
  width: 8rem;
}

.page-link {
  border: 0;
  border-radius: .45rem;
  margin-inline: .1rem;
}

@media (max-width: 575.98px) {
  .search-box {
    width: 100%;
    order: 4;
  }
}
</style>