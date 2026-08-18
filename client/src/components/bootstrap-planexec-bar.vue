<template>
  <div class="plan-progress" :class="{ 'is-compact': compact }">
    <div class="d-flex align-items-center justify-content-between gap-2 mb-1">
      <span class="small text-body-secondary">Realizado</span>
      <strong class="small" :class="statusClass">{{ percentageLabel }}</strong>
    </div>
    <div class="progress" role="progressbar"
      :aria-label="`Executado ${currencyBRL(safeExecuted)} de ${currencyBRL(safePlanned)}`"
      :aria-valuenow="Math.round(percentage)" aria-valuemin="0" aria-valuemax="100">
      <div :class="['progress-bar', barClass]" :style="{ width: visibleWidth }"></div>
      <div v-if="dividerPosition" class="divider" :style="{ left: dividerPosition }"
        title="Proporção esperada até o mês atual"></div>
    </div>
    <div v-if="!compact" class="d-flex flex-wrap align-items-center justify-content-between gap-2 mt-2 small">
      <span><strong :class="statusClass">{{ currencyBRL(safeExecuted) }}</strong> realizado</span>
      <span class="text-body-secondary">{{ differenceLabel }}</span>
    </div>
  </div>
</template>

<script setup>
import { computed } from "vue";
import { currencyBRL } from "./filters/currency.filter";

const props = defineProps({
  planned: { type: Number, default: 0 },
  executed: { type: Number, default: 0 },
  percentDivider: { type: Number, default: 0 },
  tone: { type: String, default: "primary" },
  compact: { type: Boolean, default: false },
});

const safePlanned = computed(() => Math.max(0, Number(props.planned) || 0));
const safeExecuted = computed(() => Math.max(0, Number(props.executed) || 0));
const percentage = computed(() => safePlanned.value > 0 ? (safeExecuted.value / safePlanned.value) * 100 : safeExecuted.value > 0 ? 100 : 0);
const overBudget = computed(() => safePlanned.value > 0 && safeExecuted.value > safePlanned.value);
const visibleWidth = computed(() => `${Math.min(percentage.value, 100)}%`);
const percentageLabel = computed(() => safePlanned.value === 0 && safeExecuted.value > 0 ? "Sem planejamento" : `${percentage.value.toLocaleString("pt-BR", { maximumFractionDigits: 1 })}%`);
const barClass = computed(() => overBudget.value ? "bg-danger" : `bg-${props.tone}`);
const statusClass = computed(() => overBudget.value ? "text-danger" : `text-${props.tone}`);
const dividerPosition = computed(() => props.percentDivider > 0 ? `${Math.min(Math.max(props.percentDivider, 0), 1) * 100}%` : null);
const differenceLabel = computed(() => {
  const difference = safePlanned.value - safeExecuted.value;
  if (difference > 0) return `${currencyBRL(difference)} disponível`;
  if (difference < 0) return `${currencyBRL(Math.abs(difference))} acima do planejado`;
  return safePlanned.value ? "Planejamento atingido" : "Sem valor planejado";
});
</script>

<style scoped>
.progress {
  position: relative;
  height: .6rem;
  overflow: visible;
  background: var(--bs-tertiary-bg);
  border-radius: 999px;
}

.progress-bar {
  border-radius: 999px;
  transition: width .35s ease;
}

.divider {
  position: absolute;
  top: -.2rem;
  bottom: -.2rem;
  width: 2px;
  transform: translateX(-1px);
  background: var(--bs-body-color);
  border-radius: 2px;
  opacity: .65;
}

.is-compact .progress {
  height: .45rem;
}
</style>