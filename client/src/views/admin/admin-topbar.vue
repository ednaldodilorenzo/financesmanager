<template>
  <header class="admin-topbar navbar bg-white border-bottom px-3 px-lg-4">
    <div class="d-flex align-items-center gap-3">
      <button class="btn btn-light rounded-circle topbar-button" type="button" aria-label="Alternar menu lateral"
        @click="$emit('collapse-click')">
        <i :class="['bi', sidebarCollapsed ? 'bi-layout-sidebar-inset' : 'bi-list']"></i>
      </button>
      <div class="d-none d-sm-block">
        <span class="small text-body-secondary">Você está em</span>
        <strong class="d-block lh-sm">{{ pageTitle }}</strong>
      </div>
    </div>

    <div class="d-flex align-items-center gap-1 ms-auto">
      <div class="dropdown">
        <button class="btn btn-light rounded-circle topbar-button" type="button" aria-label="Menu do usuário"
          :aria-expanded="userMenuOpen" @click.stop="toggleUserMenu"><i class="bi bi-person"></i></button>
        <ul class="dropdown-menu dropdown-menu-end border-0 shadow-sm" :class="{ show: userMenuOpen }" @click.stop>
          <li class="px-3 py-2 border-bottom">
            <strong class="small d-block">{{ userName }}</strong>
            <span class="small text-body-secondary">Minha conta</span>
          </li>
          <li><button class="dropdown-item py-2" type="button" @click="onChangePasswordClick"><i
                class="bi bi-key me-2"></i>Alterar senha</button></li>
        </ul>
      </div>
      <button class="btn btn-light rounded-circle topbar-button text-danger" type="button" title="Sair"
        aria-label="Sair" @click="logout">
        <i class="bi bi-box-arrow-right"></i>
      </button>
    </div>
  </header>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useStore } from "vuex";
import authService from "@/views/login/auth.service";
import { ROUTE_NAMES as AUTH_ROUTES } from "@/views/login/routes.definition";
import { useModalScreen } from "@/components/modal/use-modal-screen";
import ChangePassword from "../login/change-password.vue";

defineProps({ sidebarCollapsed: { type: Boolean, default: false } });
defineEmits(["collapse-click"]);

const route = useRoute();
const router = useRouter();
const store = useStore();
const modal = useModalScreen(ChangePassword);
const userMenuOpen = ref(false);
const pageTitle = computed(() => route.meta?.label || route.meta?.title || "Painel");
const currentUser = computed(() => store.getters["currentUser/getUser"] || {});
const userName = computed(() => currentUser.value.name || currentUser.value.username || "Minha conta");

const closeUserMenu = () => { userMenuOpen.value = false; };
const toggleUserMenu = () => { userMenuOpen.value = !userMenuOpen.value; };

async function onChangePasswordClick() {
  closeUserMenu();
  await modal.show();
}

async function logout() {
  await authService.logout();
  router.push({ name: AUTH_ROUTES.INDEX });
}

onMounted(() => document.addEventListener("click", closeUserMenu));
onBeforeUnmount(() => document.removeEventListener("click", closeUserMenu));
</script>

<style scoped>
.admin-topbar {
  position: sticky;
  top: 0;
  z-index: 1025;
  min-height: 4.25rem;
}

.topbar-button {
  width: 2.5rem;
  height: 2.5rem;
  display: inline-grid;
  place-items: center;
  padding: 0;
  font-size: 1.15rem;
}

.dropdown-menu.show {
  display: block;
  margin-top: .65rem;
  min-width: 13rem;
}
</style>