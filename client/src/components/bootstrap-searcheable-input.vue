<template>
  <div style="position: relative">
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
    <ul
      ref="dropdown"
      class="dropdown-menu"
      style="width: 100%"
      @click="onDropdownClick"
      v-click-outside="onDropdownClick"
    >
      <li v-for="item in options" :key="item[valueField]">
        <a class="dropdown-item" href="#" @click="onItemSelected(item)">{{
          item[displayField]
        }}</a>
      </li>
    </ul>
  </div>
</template>
<script>
export default {
  name: "BootstrapSearchableInput",
  emits: ["search-input", "update:modelValue"],
  props: {
    modelValue: {
      type: Object,
      default: () => {},
    },
    displayField: {
      type: String,
    },
    valueField: {
      type: String,
    },
    options: {
      type: Array,
      default: () => [],
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
      result: this.modelValue?.[this.displayField],
    };
  },
  methods: {
    onInput(event) {
      this.$refs.dropdown.style.display = "block";
      this.result = event.target.value;
      if (!this.result) {
        this.$emit("update:modelValue", null);
      }
      this.$emit("search-input", event);
    },
    onDropdownClick() {
      this.$refs.dropdown.style.display = "none";
    },
    onItemSelected(value) {
      this.$emit("update:modelValue", value);
    },
  },
  watch: {
    modelValue(val) {
      this.result = val?.[this.displayField];
    },
  },
};
</script>
