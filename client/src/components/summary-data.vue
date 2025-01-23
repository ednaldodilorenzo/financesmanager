<template>
  <div class="card mb-3">
    <div class="card-body p-0">
      <div
        class="row row-cols-xxl-5 row-cols-md-3 row-cols-1 g-0 text-center align-items-center justify-content-center"
      >
        <div
          v-for="item in data"
          class="col border-end border-light border-dashed"
        >
          <div class="mt-3 mt-md-0 p-3">
            <h5
              title="Number of Orders"
              class="text-muted fs-13 text-uppercase"
            >
              {{ item.title }}
            </h5>
            <div
              class="d-flex align-items-center justify-content-center gap-2 my-3"
            >
              <div class="avatar-sm flex-shrink-0">
                <span
                  class="avatar-title bg-primary-subtle fs-22 rounded-circle text-primary"
                  ><iconify-icon
                    icon="solar:bill-list-bold-duotone"
                  ></iconify-icon
                ></span>
              </div>
              <h3 class="mb-0 fw-bold">
                {{ item.value }}
              </h3>
            </div>
            <p class="mb-0 text-muted">
              <span class="text-success me-2"
                ><i class="ti ti-caret-up-filled"></i> {{ item.percent }}</span
              ><span class="text-nowrap">{{ item.percentMessage }}</span>
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup>
const props = defineProps({
  data: {
    type: Array,
    default: () => [],
  },
});
</script>
