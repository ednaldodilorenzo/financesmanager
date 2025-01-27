<template>
  <div class="container">
    <label class="form-label" v-if="label" :for="$attrs['id']">{{
      label
    }}</label>
    <div class="position-relative">
      <div
        class="form-control d-flex flex-wrap align-items-center"
        id="tag-input-wrapper"
      >
        <!-- Tags as badges inside input -->
        <span
          v-for="(tag, index) in tags"
          :key="index"
          class="badge bg-info text-white me-1 mb-1 d-flex align-items-center"
          style="gap: 8px"
        >
          {{ tag }}
          <button
            type="button"
            class="btn-close btn-close-white"
            @click.prevent.stop="removeTag(index)"
            style="font-size: 0.6rem"
            aria-label="Remove"
          ></button>
        </span>
        <!-- Input field for new tags -->
        <input
          type="text"
          v-model="newTag"
          @input="filterOptions"
          class="border-0 flex-grow-1"
          placeholder="Insira a tag"
          @keyup.enter.prevent.stop="addTag"
        />
      </div>
      <ul
        v-if="filteredOptions.length > 0"
        ref="dropdown"
        class="dropdown-menu show position-absolute w-100"
        @click="onDropdownClick"
      >
        <li
          v-for="(option, index) in filteredOptions"
          :key="index"
          :class="['dropdown-item', { active: highlightedIndex === index }]"
          @click="selectOption(option)"
          @mouseover="highlightOption(index)"
        >
          {{ option }}
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick } from "vue";
import { debounce } from "@/utils/support";

const props = defineProps({
  label: {
    type: [String, Boolean],
    default: false,
  },
  options: {
    type: [Array, Function],
    default: () => [],
  },
});

let newTag = ref("");
let tags = defineModel("modelValue");
tags.value = [];

const filteredOptions = ref([]);
const highlightedIndex = ref(-1);
const dropdown = ref(null);

const fetchOptions = (filter) => {
  if (typeof props.options === "function") {
    return props.options(filter).then((result) => result);
  } else {
    console.log(filter);

    return new Promise((resolve, reject) => {
      const result = props.options.filter((option) =>
        option.toLowerCase().includes(filter.toLowerCase())
      );

      resolve(result);
    });
  }
};

const loadItems = (value) => {
  fetchOptions(value).then((resp) => {
    filteredOptions.value = resp;
  });
};

const fetchOptionsDebounced = debounce(loadItems, 1000);

const addTag = () => {
  const tag = newTag.value.trim();
  if (tag && !tags.value.includes(tag)) {
    tags.value.push(tag);
  }
  newTag.value = "";
};

const removeTag = (index) => {
  tags.value = tags.value.toSpliced(index, 1);
};

const selectOption = async (option) => {
  newTag.value = option;
  addTag();
  // defer unting the screen is mounted.
  await nextTick();
  if (dropdown.value) {
    dropdown.value.style.display = "none";
  }
};

const filterOptions = (event) => {
  // Trigger computation of filteredOptions
  if (!event.target.value) {
    highlightedIndex.value = -1;
    filteredOptions.value = [];
  } else {
    fetchOptionsDebounced(event.target.value, 1000);
  }
};

const onDropdownClick = () => {
  if (dropdown.value) {
    dropdown.value.style.display = "none";
  }
};

const highlightOption = (index) => {
  highlightedIndex.value = index;
};
</script>
<style scoped>
#tag-input-wrapper {
  min-height: 38px; /* Ensures proper height for input */
  overflow-x: auto; /* Allows horizontal scrolling if necessary */
  padding: 0.5rem; /* Padding inside the container */
}
input::placeholder {
  color: #adb5bd; /* Placeholder color */
}

/* Remove border when input is focused */
#tag-input-wrapper input:focus {
  outline: none; /* Remove the focus outline */
  border: none; /* Remove the border */
  box-shadow: none; /* Remove any Bootstrap box-shadow for focus */
}
.dropdown-menu {
  max-height: 200px;
  overflow-y: auto;
}

.dropdown-item.active {
  background-color: #007bff;
  color: white;
}
</style>
