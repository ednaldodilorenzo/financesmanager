<template>
  <div v-if="open" class="backdrop"></div>
  <div v-if="open" class="modal" style="display: block" tabindex="-1">
    <div class="modal-dialog modal-lg">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">{{ title }}</h5>
          <button
            type="button"
            class="btn-close"
            @click="closeModal()"
            aria-label="Close"
          ></button>
        </div>
        <div class="modal-body">
          <slot></slot>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" @click="closeModal()">
            Fechar
          </button>
          <button type="button" class="btn btn-primary" @click="saveClick()">
            Salvar
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
export default {
  name: "BoostrapModal",
  emits: ["close", "save"],
  props: {
    visible: {
      type: Boolean,
      default: false,
    },
    title: {
      type: String,
      default: "",
    },
  },
  data() {
    return {
      open: this.visible,
    };
  },
  methods: {
    closeModal() {
      this.$emit("close");
    },
    saveClick() {
      this.$emit("save");
    },
  },
  watch: {
    visible: function (newVal, oldVal) {
      this.open = newVal;
    },
  },
};
</script>
