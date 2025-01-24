<template>
  <tr>
    <td>
      <span v-if="!editMode">{{ item?.category?.name }}</span>
      <bootstrap-searcheable-select
        display-field="name"
        value-field="id"
        v-model="item.category"
        :options="options"
        autocomplete="off"
        v-else
      />
    </td>
    <td>
      <span v-if="!editMode">{{ item.strValue }}</span>
      <input v-model="item.strValue" class="form-control" v-currency v-else />
    </td>
    <td>
      <button
        class="btn icon-link link-primary icon-link-hover"
        v-if="!editMode"
        @click="onEditClick"
      >
        <i class="bi bi-pencil-fill"></i>
      </button>
      <button
        class="btn icon-link link-danger icon-link-hover"
        v-if="!editMode"
        @click="onDeleteClick"
      >
        <i class="bi bi-trash-fill"></i>
      </button>
      <button
        class="btn icon-link link-success icon-link-hover"
        v-if="editMode"
        @click="onSaveClick"
      >
        <i class="bi bi-check-lg"></i>
      </button>
      <button
        v-if="editMode"
        @click="onCancel"
        class="btn icon-link link-danger icon-link-hover"
      >
        <i class="bi bi-x-lg"></i>
      </button>
    </td>
  </tr>
</template>
<script setup>
import { ref } from "vue";
import { parseCurrencyToNumber } from "@/utils/numbers";
import BootstrapInput from "@/components/bootstrap-input.vue";
import BootstrapSearcheableSelect from "@/components/bootstrap-searcheable-select.vue";
import budgetService from "./budget.service";
import { useToast } from "vue-toastification";
import { useLoadingScreen } from "@/components/loading/useLoadingScreen";
import { useDialogScreen } from "@/components/dialog/use-dialog-screen";

const toast = useToast();
const loading = useLoadingScreen();
const dialog = useDialogScreen(
  "Deseja excluir categoria do orçamento?",
  "Exclusão"
);

const props = defineProps({
  options: { type: Array, default: () => [] },
  budgetItem: { type: Object, default: null },
});

let editing = false;

const emit = defineEmits(["cancel-item", "delete-item", "change-item"]);

let editMode = props.budgetItem.id ? ref(false) : ref(true);
let item = ref(props.budgetItem.clone());

const onCancel = () => {
  if (editing) {
    editing = false;
    item.value = props.budgetItem.clone();
    editMode.value = false;
  } else {
    emit("cancel-item", props.budgetItem);
  }
};

const onDeleteClick = async () => {
  const confirm = await dialog.show();
  if (confirm) {
    loading.show();
    budgetService
      .delete(item.value.id)
      .then(() => {
        toast.success("Categoria excluída com sucesso!", {
          position: "top-center",
        });
        emit("delete-item", item.value);
      })
      .finally(() => {
        loading.hide();
      });
  }
};

const onEditClick = () => {
  editMode.value = true;
  editing = true;
};

const onSaveClick = () => {
  loading.show();
  // const { category, ...payload } = {
  //   ...item.value,
  //   categoryId: item.value.category.id,
  //   value: parseCurrencyToNumber(item.value.value),
  // };

  const method = item.value.id
    ? budgetService.modify(item.value.id, item.value)
    : budgetService.create(item.value);

  method
    .then(() => {
      editMode.value = false;
      editing = true;
      toast.success(
        `Categoria ${item.value.id ? "atualizada" : "incluída"} com sucesso!`,
        {
          position: "top-center",
        }
      );
      emit("change-item", item.value.clone());
    })
    .finally(() => {
      loading.hide();
    });
};
</script>
