<template>
  <div class="progress" style="height: 20px">
    <div
      class="progress-bar"
      role="progressbar"
      :style="{
        width: executed <= planned ? (executed / planned) * 100 + '%' : '100%',
      }"
      aria-valuenow="25"
      aria-valuemin="0"
      aria-valuemax="100"
    >
      {{ currencyBRL(executed) }}
    </div>
    <div
      v-if="executed > planned"
      class="progress-bar bg-danger"
      role="progressbar"
      :style="{
        width: (executed / planned - 1) * 100 + '%',
      }"
      aria-valuenow="30"
      aria-valuemin="0"
      aria-valuemax="100"
    >
      {{ currencyBRL(executed - planned) }}
    </div>
    <div class="flex-fill text-center" v-else>
      {{ currencyBRL(planned - executed) }}
    </div>
  </div>
</template>
<script setup>
import { currencyBRL } from "./filters/currency.filter";

const props = defineProps({
  planned: {
    type: Number,
    default: () => 0,
  },
  executed: {
    type: Number,
    default: () => 0,
  },
});
</script>
