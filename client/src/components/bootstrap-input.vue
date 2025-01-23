<template>
  <div class="container">
    <label class="form-label" v-if="label" :for="$attrs['id']">{{
      label
    }}</label>
    <input
      v-bind="$attrs"
      class="form-control"
      :value="result"
      @input="onInput"
      :required="required"
    />
    <div class="invalid-feedback" id="live-feedback-email">
      {{ requiredMessage }}
    </div>
  </div>
</template>
<script>
export default {
  name: "BootstrapInput",
  emits: ["update:modelValue"],
  props: {
    modelValue: {
      type: String,
      default: "",
    },
    required: {
      type: Boolean,
      default: false,
    },
    requiredMessage: {
      type: String,
      default: "Campo obrigatório",
    },
    label: {
      type: [String, Boolean],
      default: false,
    },
  },
  data() {
    return {
      result: this.modelValue,
    };
  },
  methods: {
    onInput(event) {
      this.result = event.target.value;
      this.$emit("update:modelValue", this.result);
    },
  },
  watch: {
    modelValue: function (newVal, oldVal) {
      this.result = newVal;
    },
  },
};
</script>
