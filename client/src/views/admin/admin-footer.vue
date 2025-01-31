<template>
  <div
    class="navbar d-flex justify-content-end bg-body-tertiary"
    style="height: 50px; margin-top: auto"
  >
    <span class="mx-3">Versão {{ appVersion }}</span> 
  </div>
</template>
<script>
export default {
  name: "admin-footer",
  data() {
    return {
      appVersion: process.env.VUE_APP_VERSION,
    };
  },
};
</script>
