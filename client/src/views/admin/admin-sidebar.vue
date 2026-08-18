<template>
  <aside class="admin-sidebar d-flex flex-column"
    :class="{ collapsed: collapsed && !isMobile, 'mobile-open': mobileOpen }" :aria-hidden="isMobile && !mobileOpen">
    <div class="brand d-flex align-items-center gap-3 px-3">
      <span class="brand-mark"><i class="bi bi-wallet2"></i></span>
      <span class="brand-name fw-bold text-nowrap">Meus Serviços</span>
    </div>

    <div class="sidebar-divider mx-3"></div>
    <span class="menu-label px-3 mb-2">MENU</span>

    <nav class="flex-grow-1 overflow-y-auto px-2" aria-label="Menu principal">
      <ul class="nav nav-pills flex-column gap-1">
        <li v-for="menuItem in menuItems" :key="menuItem.name" class="nav-item">
          <router-link :to="menuItem" class="nav-link d-flex align-items-center gap-3"
            :title="collapsed && !isMobile ? menuItem.meta.label : undefined" @click="$emit('navigate')">
            <i :class="[menuItem.meta.icon, 'menu-icon bx']"></i>
            <span class="menu-text text-nowrap">{{ menuItem.meta.label }}</span>
          </router-link>
        </li>
      </ul>
    </nav>

    <div class="user-panel d-flex align-items-center gap-3 m-2 p-2">
      <span class="user-avatar"><i class="bi bi-person"></i></span>
      <div class="menu-text min-w-0">
        <strong class="small d-block text-truncate">{{ userName }}</strong>
        <span class="user-role d-block text-truncate">Usuário</span>
      </div>
    </div>
  </aside>
</template>

<script setup>
import { computed, shallowRef } from "vue";
import { useStore } from "vuex";
import { SIDEBAR_ROUTES } from "@/views/routes.definition";

defineProps({
  collapsed: { type: Boolean, default: false },
  mobileOpen: { type: Boolean, default: false },
  isMobile: { type: Boolean, default: false },
});
defineEmits(["navigate"]);

const store = useStore();
const menuItems = shallowRef(SIDEBAR_ROUTES);
const currentUser = computed(() => store.getters["currentUser/getUser"] || {});
const userName = computed(() => currentUser.value.name || currentUser.value.username || "Minha conta");
</script>

<style scoped>
.admin-sidebar {
  position: sticky;
  top: 0;
  z-index: 1040;
  width: 16.5rem;
  height: 100vh;
  flex: 0 0 auto;
  color: #fff;
  background: linear-gradient(180deg, #4144db 0%, #3437bd 100%);
  transition: width .22s ease, transform .22s ease;
  box-shadow: .25rem 0 1.5rem rgba(31, 38, 135, .14);
}

.brand {
  height: 4.25rem;
  overflow: hidden;
}

.brand-mark {
  width: 2.5rem;
  height: 2.5rem;
  display: grid;
  place-items: center;
  flex: 0 0 auto;
  border-radius: .75rem;
  background: rgba(255, 255, 255, .14);
  font-size: 1.2rem;
}

.brand-name {
  font-size: 1.05rem;
}

.sidebar-divider {
  height: 1px;
  background: rgba(255, 255, 255, .13);
}

.menu-label {
  margin-top: 1.15rem;
  color: #bfc6ff;
  font-size: .67rem;
  font-weight: 700;
  letter-spacing: .12em;
  transition: opacity .15s ease;
}

.nav-link {
  min-height: 2.75rem;
  padding: .65rem .75rem;
  overflow: hidden;
  color: #dfe2ff;
  border-radius: .65rem;
  transition: color .15s ease, background-color .15s ease;
}

.nav-link:hover {
  color: #fff;
  background: rgba(255, 255, 255, .1);
}

.nav-link.router-link-active {
  color: #fff;
  background: rgba(255, 255, 255, .16);
  box-shadow: inset 3px 0 0 #fff;
}

.menu-icon {
  width: 1.5rem;
  flex: 0 0 auto;
  text-align: center;
  font-size: 1.25rem;
}

.menu-text {
  opacity: 1;
  transition: opacity .12s ease;
}

.user-panel {
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, .1);
  border-radius: .75rem;
  background: rgba(255, 255, 255, .07);
}

.user-avatar {
  width: 2.25rem;
  height: 2.25rem;
  display: grid;
  place-items: center;
  flex: 0 0 auto;
  border-radius: 50%;
  background: rgba(255, 255, 255, .14);
}

.user-role {
  color: #c7ccf7;
  font-size: .72rem;
}

.min-w-0 {
  min-width: 0;
}

.admin-sidebar.collapsed {
  width: 4.75rem;
}

.admin-sidebar.collapsed .brand-name,
.admin-sidebar.collapsed .menu-text,
.admin-sidebar.collapsed .menu-label {
  width: 0;
  opacity: 0;
  pointer-events: none;
}

.admin-sidebar.collapsed .brand,
.admin-sidebar.collapsed .nav-link {
  justify-content: center;
}

.admin-sidebar.collapsed .brand {
  padding-inline: 1.1rem !important;
}

.admin-sidebar.collapsed .user-panel {
  justify-content: center;
}

@media (max-width: 991.98px) {
  .admin-sidebar {
    position: fixed;
    left: 0;
    transform: translateX(-105%);
    width: min(17rem, 86vw);
  }

  .admin-sidebar.mobile-open {
    transform: translateX(0);
  }
}
</style>