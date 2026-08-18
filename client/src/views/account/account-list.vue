<template>
  <section class="accounts-page py-3">
    <header class="d-flex flex-wrap align-items-center justify-content-between gap-3 mb-4">
      <div>
        <h1 class="h4 fw-bold mb-1">Contas</h1>
        <p class="small text-body-secondary mb-0">Gerencie bancos, cartões, dinheiro e investimentos</p>
      </div>
      <button class="btn btn-primary px-3" type="button" @click="onNewClicked">
        <i class="bi bi-plus-lg me-1"></i>Nova conta
      </button>
    </header>

    <div class="card border-0 shadow-sm mb-4">
      <div class="card-body p-3">
        <div class="d-flex flex-wrap align-items-center justify-content-between gap-3">
          <div>
            <span class="small text-body-secondary d-block">Contas cadastradas</span>
            <strong class="h5 mb-0">{{ items.length }}</strong>
          </div>
          <div class="search-box">
            <i class="bi bi-search"></i>
            <input v-model="search" class="form-control" type="search" placeholder="Pesquisar conta…"
              aria-label="Pesquisar conta" />
          </div>
        </div>
      </div>
    </div>

    <div v-if="filteredItems.length" class="row g-3">
      <div class="col-sm-6 col-xl-4 col-xxl-3">
        <button class="new-account-card card border-0 shadow-sm h-100 w-100 text-start" type="button"
          @click="onNewClicked">
          <div class="card-body d-flex align-items-center gap-3 p-4">
            <span class="new-icon"><i class="bi bi-plus-lg"></i></span>
            <div>
              <strong class="d-block">Adicionar conta</strong>
              <span class="small text-body-secondary">Cadastre uma nova conta</span>
            </div>
          </div>
        </button>
      </div>
      <div v-for="item in filteredItems" :key="item.id" class="col-sm-6 col-xl-4 col-xxl-3">
        <AccountItem :item="item" @item-edit-click="onItemEditClick" />
      </div>
    </div>

    <div v-else class="card border-0 shadow-sm">
      <div class="card-body text-center py-5 px-3">
        <span class="empty-icon"><i class="bi bi-wallet2"></i></span>
        <h2 class="h6 mt-3 mb-1">{{ search ? 'Nenhuma conta encontrada' : 'Nenhuma conta cadastrada' }}</h2>
        <p class="small text-body-secondary mb-3">{{ search ? 'Tente pesquisar por outro nome.' : 'Cadastre sua primeira conta para começar.' }}</p>
        <button v-if="!search" class="btn btn-primary btn-sm" type="button" @click="onNewClicked"><i
            class="bi bi-plus-lg me-1"></i>Nova conta</button>
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed, ref } from "vue";
import { useRouter } from "vue-router";
import accountService from "./account.service";
import { useLoadingScreen } from "@/components/loading/useLoadingScreen";
import { useModalScreen } from "@/components/modal/use-modal-screen";
import AccountChangeScreen from "./account-change-screen.vue";
import AccountItem from "./account-item.vue";

const loading = useLoadingScreen();
const router = useRouter();
const modal = useModalScreen(AccountChangeScreen);
const items = ref([]);
const search = ref("");

const filteredItems = computed(() => {
  const term = search.value.trim().toLocaleLowerCase("pt-BR");
  return term ? items.value.filter(({ name }) => name.toLocaleLowerCase("pt-BR").includes(term)) : items.value;
});

async function getList() {
  loading.show();
  try {
    const response = await accountService.findAll({ paginate: false });
    items.value = response.data;
  } catch (error) {
    router.push({ name: "denied" });
  } finally { loading.hide(); }
}

async function onItemEditClick(item) {
  if (await modal.show(item)) getList();
}

async function onNewClicked() {
  if (await modal.show()) getList();
}

getList();
</script>

<style scoped>
.new-account-card {
  min-height: 8.5rem;
  color: var(--bs-body-color);
  background: rgba(var(--bs-primary-rgb), .035);
  border: 1px dashed rgba(var(--bs-primary-rgb), .35) !important;
  transition: transform .15s ease, background-color .15s ease;
}

.new-account-card:hover {
  transform: translateY(-2px);
  background: rgba(var(--bs-primary-rgb), .07);
}

.new-icon,
.empty-icon {
  display: grid;
  place-items: center;
  background: rgba(var(--bs-primary-rgb), .1);
  color: var(--bs-primary);
}

.new-icon {
  width: 2.75rem;
  height: 2.75rem;
  flex: 0 0 auto;
  border-radius: .8rem;
}

.empty-icon {
  width: 3.25rem;
  height: 3.25rem;
  display: inline-grid;
  border-radius: 1rem;
  font-size: 1.35rem;
}
</style>