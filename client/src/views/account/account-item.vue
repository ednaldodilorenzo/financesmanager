<template>
  <article class="account-card card border-0 shadow-sm h-100">
    <div class="card-body p-4">
      <div class="d-flex align-items-start justify-content-between gap-3">
        <span :class="['account-icon', typeInfo.tone]"><i :class="['bi', typeInfo.icon]"></i></span>
        <div class="dropdown">
          <button class="btn btn-sm btn-light rounded-circle" type="button" aria-label="Ações da conta"
            :aria-expanded="menuOpen" @click.stop="toggleAccountMenu">
            <i class="bi bi-three-dots-vertical"></i>
          </button>
          <ul class="dropdown-menu dropdown-menu-end border-0 shadow-sm" :class="{ show: menuOpen }" @click.stop>
            <li><button class="dropdown-item" type="button" @click="onEditClick"><i class="bi bi-pencil me-2"></i>Editar
                conta</button></li>
          </ul>
        </div>
      </div>
      <h2 class="h6 fw-bold text-truncate mt-3 mb-1" :title="item.name">{{ item.name }}</h2>
      <span class="small text-body-secondary">{{ typeInfo.label }}</span>
      <div v-if="item.type === 'C' && item.dueDay" class="small mt-3 pt-3 border-top text-body-secondary">
        <i class="bi bi-calendar3 me-1"></i>Pagamento no dia <strong class="text-body">{{ item.dueDay }}</strong>
      </div>
    </div>
  </article>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from "vue";

const props = defineProps({ item: { type: Object, required: true } });
const emit = defineEmits(["item-edit-click"]);
const menuOpen = ref(false);

const accountTypes = {
  A: { label: "Conta corrente", icon: "bi-bank", tone: "is-current" },
  C: { label: "Cartão de crédito", icon: "bi-credit-card", tone: "is-credit" },
  D: { label: "Dinheiro", icon: "bi-cash-stack", tone: "is-cash" },
  I: { label: "Investimento", icon: "bi-graph-up-arrow", tone: "is-investment" },
};
const typeInfo = computed(() => accountTypes[props.item.type] || { label: "Conta", icon: "bi-wallet2", tone: "" });
const closeMenu = () => { menuOpen.value = false; };
const toggleAccountMenu = () => { menuOpen.value = !menuOpen.value; };

function onEditClick() {
  closeMenu();
  emit("item-edit-click", props.item);
}

onMounted(() => document.addEventListener("click", closeMenu));
onBeforeUnmount(() => document.removeEventListener("click", closeMenu));
</script>

<style scoped>
.account-card {
  min-height: 8.5rem;
  border-radius: .85rem;
  transition: transform .15s ease, box-shadow .15s ease;
}

.account-card:hover {
  transform: translateY(-2px);
}

.account-icon {
  width: 2.75rem;
  height: 2.75rem;
  display: grid;
  place-items: center;
  border-radius: .8rem;
  background: var(--bs-tertiary-bg);
  color: var(--bs-secondary);
  font-size: 1.15rem;
}

.account-icon.is-current {
  color: var(--bs-primary);
  background: rgba(var(--bs-primary-rgb), .1);
}

.account-icon.is-credit {
  color: #7c3aed;
  background: rgba(124, 58, 237, .1);
}

.account-icon.is-cash {
  color: var(--bs-success);
  background: rgba(var(--bs-success-rgb), .1);
}

.account-icon.is-investment {
  color: #d97706;
  background: rgba(217, 119, 6, .1);
}

.dropdown-menu.show {
  display: block;
  z-index: 1080;
}
</style>