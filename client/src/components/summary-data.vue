<template>
  <div :class="layout === 'sidebar' ? 'd-grid gap-3' : 'row row-cols-1 row-cols-md-2 row-cols-xxl-4 g-3'">
    <div v-for="item in data" :key="item.title" :class="layout === 'sidebar' ? '' : 'col'">
      <div class="summary-card card border-0 h-100 shadow-sm">
        <div class="card-body p-3 p-lg-4">
          <div class="d-flex align-items-start justify-content-between gap-3">
            <div>
              <p class="small fw-semibold text-body-secondary mb-2">{{ item.title }}</p>
              <h3 class="h4 fw-bold mb-0" :class="valueClass(item)">{{ $filters.currencyBRL(item.value || 0) }}</h3>
            </div>
            <span :class="['summary-icon', `text-bg-${item.tone || 'primary'}`]">
              <i :class="['bi', item.icon || 'bi-wallet2']"></i>
            </span>
          </div>
          <p v-if="item.percentMessage" class="small text-body-secondary mt-3 mb-0">
            <span class="fw-semibold me-1">{{ $filters.percentageBRL(item.percent || 0) }}</span>{{ item.percentMessage
            }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  data: { type: Array, default: () => [] },
  layout: { type: String, default: "grid", validator: (value) => ["grid", "sidebar"].includes(value) },
});
const valueClass = (item) => item.value < 0 ? "text-danger" : item.tone ? `text-${item.tone}` : "text-body";
</script>

<style scoped>
.summary-card {
  background: var(--bs-tertiary-bg);
  transition: transform .15s ease, box-shadow .15s ease;
}

.summary-card:hover {
  transform: translateY(-2px);
}

.summary-icon {
  width: 2.5rem;
  height: 2.5rem;
  display: grid;
  place-items: center;
  border-radius: .75rem;
  --bs-bg-opacity: .12;
}
</style>