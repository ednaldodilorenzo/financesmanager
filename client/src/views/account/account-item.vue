<template>
  <div class="card" style="width: 18rem; margin: 0.5rem">
    <div class="card-body">
      <div class="d-flex justify-content-between">
        <h5 class="card-title">{{ props.item.name }}</h5>
        <button
          class="btn"
          :class="{ show: accountMenuToogle }"
          type="button"
          id="defaultDropdown"
          data-bs-toggle="dropdown"
          data-bs-auto-close="true"
          aria-expanded="false"
          @click="toggleAccountMenu"
        >
          <i class="bi bi-three-dots-vertical"></i>
        </button>

        <ul
          class="dropdown-menu"
          :data-bs-popper="accountMenuToogle ? 'static' : null"
          :class="{ show: accountMenuToogle }"
          aria-labelledby="defaultDropdown"
        >
          <li>
            <a class="dropdown-item" href="#" @click.prevent="onEditClick(props.item)"
              >Editar</a
            >
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>
<script setup>
import { ref } from "vue";

const accountMenuToogle = ref(false);
const emit = defineEmits(["item-edit-click"]);

const props = defineProps({
  item: {
    type: Object,
    default: () => {},
  },
});

const toggleAccountMenu = () => {
  accountMenuToogle.value = !accountMenuToogle.value;
};

const onEditClick = (item) => {
  accountMenuToogle.value = false;
  emit("item-edit-click", item);
};
</script>
