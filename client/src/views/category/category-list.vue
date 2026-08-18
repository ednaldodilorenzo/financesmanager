<template>
  <section class="categories-page py-3">
    <header class="d-flex flex-wrap align-items-center justify-content-between gap-3 mb-4">
      <div>
        <h1 class="h4 fw-bold mb-1">Categorias</h1>
        <p class="small text-body-secondary mb-0">Organize receitas, despesas e investimentos</p>
      </div>
      <button class="btn btn-primary px-3" type="button" @click="onNewClicked"><i class="bi bi-plus-lg me-1"></i>Nova
        categoria</button>
    </header>

    <div class="row g-3 mb-4">
      <div v-for="summary in summaries" :key="summary.type" class="col-sm-6 col-xl-3">
        <button class="summary-card card border-0 shadow-sm w-100 h-100 text-start"
          :class="{ active: typeFilter === summary.type }" type="button" @click="toggleType(summary.type)">
          <div class="card-body d-flex align-items-center gap-3 p-3">
            <span :class="['summary-icon', summary.tone]"><i :class="['bi', summary.icon]"></i></span>
            <div><strong class="h5 d-block mb-0">{{ summary.count }}</strong><span class="small text-body-secondary">{{
                summary.label }}</span></div>
          </div>
        </button>
      </div>
      <div class="col-sm-6 col-xl-3">
        <button class="new-category-card card border-0 shadow-sm w-100 h-100 text-start" type="button"
          @click="onNewClicked">
          <div class="card-body d-flex align-items-center gap-3 p-3"><span class="new-icon"><i
                class="bi bi-plus-lg"></i></span>
            <div><strong class="d-block">Adicionar</strong><span class="small text-body-secondary">Nova categoria</span>
            </div>
          </div>
        </button>
      </div>
    </div>

    <div class="card border-0 shadow-sm categories-card">
      <div class="card-header bg-white border-bottom p-3">
        <div class="d-flex flex-wrap align-items-center justify-content-between gap-3">
          <div>
            <h2 class="h6 fw-bold mb-1">Categorias cadastradas</h2><span class="small text-body-secondary">{{
              filteredItems.length }} resultados</span>
          </div>
          <div class="search-box"><i class="bi bi-search"></i><input v-model="search" class="form-control" type="search"
              placeholder="Pesquisar categoria…" aria-label="Pesquisar categoria" /></div>
        </div>
      </div>
      <div v-if="filteredItems.length" class="table-responsive">
        <table class="table table-hover align-middle mb-0">
          <thead>
            <tr>
              <th>Categoria</th>
              <th>Tipo</th>
              <th class="text-center actions-column">Ações</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in filteredItems" :key="item.id">
              <td>
                <div class="d-flex align-items-center gap-3"><span
                    :class="['category-icon', typeInfo(item.type).tone]"><i
                      :class="['bi', typeInfo(item.type).icon]"></i></span><strong>{{ item.name }}</strong></div>
              </td>
              <td><span :class="['badge rounded-pill', typeInfo(item.type).badge]">{{ typeInfo(item.type).label
                  }}</span></td>
              <td class="text-center"><button class="btn btn-sm btn-light rounded-circle text-primary" type="button"
                  title="Editar categoria" @click="handleEdit(item)"><i class="bi bi-pencil"></i></button></td>
            </tr>
          </tbody>
        </table>
      </div>
      <div v-else class="text-center py-5 px-3">
        <span class="empty-icon"><i class="bi bi-tags"></i></span>
        <h3 class="h6 mt-3 mb-1">Nenhuma categoria encontrada</h3>
        <p class="small text-body-secondary mb-3">Altere os filtros ou cadastre uma nova categoria.</p>
        <button class="btn btn-primary btn-sm" type="button" @click="onNewClicked"><i
            class="bi bi-plus-lg me-1"></i>Nova categoria</button>
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed, ref } from "vue";
import { useRouter } from "vue-router";
import categoryService from "./category.service";
import CategoryChange from "./category-change.vue";
import { useLoadingScreen } from "@/components/loading/useLoadingScreen";
import { useModalScreen } from "@/components/modal/use-modal-screen";

const loading = useLoadingScreen();
const router = useRouter();
const modal = useModalScreen(CategoryChange);
const items = ref([]);
const search = ref("");
const typeFilter = ref("");
const types = {
  R: { label: "Receita", plural: "Receitas", icon: "bi-arrow-down-left", tone: "is-income", badge: "text-bg-success" },
  D: { label: "Despesa", plural: "Despesas", icon: "bi-arrow-up-right", tone: "is-expense", badge: "text-bg-danger" },
  I: { label: "Investimento", plural: "Investimentos", icon: "bi-graph-up-arrow", tone: "is-investment", badge: "text-bg-primary" },
};
const typeInfo = (type) => types[type] || { label: "Categoria", icon: "bi-tag", tone: "", badge: "text-bg-secondary" };
const summaries = computed(() => Object.entries(types).map(([type, info]) => ({ type, label: info.plural, icon: info.icon, tone: info.tone, count: items.value.filter((item) => item.type === type).length })));
const filteredItems = computed(() => {
  const term = search.value.trim().toLocaleLowerCase("pt-BR");
  return items.value.filter((item) => (!typeFilter.value || item.type === typeFilter.value) && (!term || item.name.toLocaleLowerCase("pt-BR").includes(term)));
});
const toggleType = (type) => { typeFilter.value = typeFilter.value === type ? "" : type; };

async function getList() {
  loading.show();
  try { items.value = (await categoryService.findAll({ paginate: false })).data; }
  catch (error) { router.push({ name: "denied" }); }
  finally { loading.hide(); }
}
async function handleEdit(item) { if (await modal.show(item)) getList(); }
async function onNewClicked() { if (await modal.show()) getList(); }
getList();
</script>

<style scoped>
.summary-card,
.new-category-card,
.categories-card {
  border-radius: .85rem
}

.summary-card {
  color: var(--bs-body-color);
  transition: transform .15s ease, box-shadow .15s ease
}

.summary-card:hover,
.summary-card.active {
  transform: translateY(-2px);
  box-shadow: 0 .5rem 1rem rgba(0, 0, 0, .1) !important
}

.summary-card.active {
  outline: 2px solid rgba(var(--bs-primary-rgb), .35)
}

.summary-icon,
.new-icon,
.category-icon,
.empty-icon {
  display: grid;
  place-items: center;
  flex: 0 0 auto
}

.summary-icon,
.new-icon,
.category-icon {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: .75rem
}

.new-category-card {
  color: var(--bs-body-color);
  background: rgba(var(--bs-primary-rgb), .035);
  border: 1px dashed rgba(var(--bs-primary-rgb), .35) !important
}

.new-icon {
  color: var(--bs-primary);
  background: rgba(var(--bs-primary-rgb), .1)
}

.categories-card .card-header {
  border-radius: .85rem .85rem 0 0
}

.actions-column {
  width: 7rem
}

.empty-icon {
  width: 3.25rem;
  height: 3.25rem;
  display: inline-grid;
  border-radius: 1rem;
  color: var(--bs-primary);
  background: rgba(var(--bs-primary-rgb), .1);
  font-size: 1.35rem
}
</style>