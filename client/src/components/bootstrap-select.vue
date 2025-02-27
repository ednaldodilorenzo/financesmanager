<template>
  <div class="container">
    <label class="form-label" v-if="label" :for="$attrs['id']">{{
      label
    }}</label>
    <select
      @change="onSelectChange"
      v-bind="$attrs"
      v-model="model"
      class="form-select"
    >
      <option v-for="option in options" :value="option[props.keyField]">
        {{ option[props.valueField] }}
      </option>
    </select>
    <div class="invalid-feedback" id="live-feedback-email">
      {{ requiredMessage }}
    </div>
  </div>
</template>
<script setup>
const props = defineProps({
  options: {
    type: Array,
    default: () => [],
  },
  keyField: {
    type: String,
    default: "",
  },
  valueField: {
    type: String,
    default: "",
  },
  label: {
    type: String,
    default: "",
  },
  requiredMessage: {
    type: String,
    default: "Campo obrigatório",
  },
});
// The name value of the parameter must be modelValue to avoid passing it externally.
const model = defineModel("modelValue");

const emit = defineEmits(["change"]);

const onSelectChange = (e) => {
  emit("change", e.target.value);
};
</script>
