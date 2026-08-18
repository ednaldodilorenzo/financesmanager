<template>
  <footer class="admin-footer bg-white border-top px-3 px-lg-4 py-3">
    <div class="d-flex flex-wrap align-items-center justify-content-between gap-2 small text-body-secondary">
      <span>&copy; {{ currentYear }} Meus Serviços</span>
      <span>Versão {{ appVersion }}</span>
    </div>
  </footer>
</template>

<script setup>
const currentYear = new Date().getFullYear();
const appVersion = process.env.VUE_APP_VERSION || "—";
</script>

<style scoped>
.admin-footer {
  margin-top: auto;
}
</style>