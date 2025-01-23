<template>
  <div v-if="loading" class="backdrop">
    <div class="spinner-border" role="status"></div>
  </div>
</template>
<script>
export default {
  props: {
    loading: {
      type: Boolean,
      required: false,
      default: false,
    },
  },
};
</script>
