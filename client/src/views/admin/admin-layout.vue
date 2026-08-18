<template>
  <div class="admin-layout d-flex min-vh-100">
    <AdminSidebar :collapsed="sidebarCollapsed" :mobile-open="mobileOpen" :is-mobile="isMobile"
      @navigate="closeMobileMenu" />

    <button v-if="isMobile && mobileOpen" class="sidebar-backdrop" type="button" aria-label="Fechar menu"
      @click="closeMobileMenu"></button>

    <div class="admin-shell d-flex flex-column flex-grow-1 min-vw-0">
      <AdminTopbar :sidebar-collapsed="sidebarCollapsed" @collapse-click="toggleSidebar" />
      <main class="admin-content flex-grow-1">
        <div class="container-fluid p-3 p-lg-4">
          <router-view />
        </div>
      </main>
      <AdminFooter />
    </div>
  </div>
</template>

<script setup>
import { onBeforeUnmount, onMounted, ref } from "vue";
import AdminFooter from "./admin-footer.vue";
import AdminTopbar from "./admin-topbar.vue";
import AdminSidebar from "./admin-sidebar.vue";

const BREAKPOINT = 992;
const isMobile = ref(window.innerWidth < BREAKPOINT);
const mobileOpen = ref(false);
const sidebarCollapsed = ref(localStorage.getItem("sidebar-collapsed") === "true");

function toggleSidebar() {
  if (isMobile.value) {
    mobileOpen.value = !mobileOpen.value;
    return;
  }
  sidebarCollapsed.value = !sidebarCollapsed.value;
  localStorage.setItem("sidebar-collapsed", String(sidebarCollapsed.value));
}

function closeMobileMenu() {
  mobileOpen.value = false;
}

function handleResize() {
  const mobile = window.innerWidth < BREAKPOINT;
  if (mobile !== isMobile.value) {
    isMobile.value = mobile;
    mobileOpen.value = false;
  }
}

onMounted(() => window.addEventListener("resize", handleResize, { passive: true }));
onBeforeUnmount(() => window.removeEventListener("resize", handleResize));
</script>

<style scoped>
.admin-layout {
  background: var(--bs-tertiary-bg);
}

.admin-shell {
  min-width: 0;
}

.admin-content {
  background: #f5f7f9;
}

.sidebar-backdrop {
  position: fixed;
  inset: 0;
  z-index: 1035;
  border: 0;
  background: rgba(15, 23, 42, .48);
  backdrop-filter: blur(2px);
}

.min-vw-0 {
  min-width: 0;
}
</style>