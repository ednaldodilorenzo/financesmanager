<template>
  <tr :class="{ 'editing-row': editMode }">
    <td>
      <div v-if="!editMode" class="d-flex align-items-center gap-3">
        <span :class="['category-icon', categoryTone]">
          <i :class="categoryIcon"></i>
        </span>
        <div class="min-w-0">
          <div class="fw-semibold text-truncate">{{ item.category?.name }}</div>
          <div class="small text-body-secondary">{{ categoryLabel }}</div>
        </div>
      </div>
      <BootstrapSearcheableSelect v-else display-field="name" value-field="id" v-model="item.category"
        :options="options" autocomplete="off" />
    </td>

    <td class="text-end value-column">
      <strong v-if="!editMode" :class="valueClass">{{ item.strValue }}</strong>
      <div v-else class="input-group input-group-sm ms-auto value-input">
        <span class="input-group-text">R$</span>
        <input v-model="item.strValue" class="form-control text-end" v-currency aria-label="Valor planejado" />
      </div>
    </td>

    <td class="text-center text-nowrap">
      <template v-if="!editMode">
        <button class="btn btn-sm btn-light rounded-circle text-primary me-1" type="button" title="Editar"
          @click="onEditClick">
          <i class="bi bi-pencil"></i>
        </button>
        <button class="btn btn-sm btn-light rounded-circle text-danger" type="button" title="Excluir"
          @click="onDeleteClick">
          <i class="bi bi-trash"></i>
        </button>
      </template>
      <template v-else>
        <button class="btn btn-sm btn-success me-1" type="button" @click="onSaveClick">
          <i class="bi bi-check-lg me-1"></i>Salvar
        </button>
        <button class="btn btn-sm btn-outline-secondary" type="button" @click="onCancel" aria-label="Cancelar">
          <i class="bi bi-x-lg"></i>
        </button>
      </template>
    </td>
  </tr>
</template>

<script setup>
import { computed, ref, watch } from "vue";
import BootstrapSearcheableSelect from "@/components/bootstrap-searcheable-select.vue";
import budgetService from "./budget.service";
import { useToast } from "vue-toastification";
import { useLoadingScreen } from "@/components/loading/useLoadingScreen";
import { useDialogScreen } from "@/components/dialog/use-dialog-screen";

const props = defineProps({
  options: { type: Array, default: () => [] },
  budgetItem: { type: Object, required: true },
});

const emit = defineEmits(["cancel-item", "delete-item", "change-item"]);
const toast = useToast();
const loading = useLoadingScreen();
const dialog = useDialogScreen("Deseja excluir a categoria do orçamento?", "Exclusão");
const editMode = ref(!props.budgetItem.id);
const editingExisting = ref(false);
const item = ref(props.budgetItem.clone());

watch(() => props.budgetItem, (value) => {
  item.value = value.clone();
  editMode.value = !value.id;
  editingExisting.value = false;
});

const categoryTone = computed(() => ({ R: "is-income", D: "is-expense", I: "is-investment" }[item.value.category?.type] || ""));
const categoryIcon = computed(() => ({ R: "bi bi-arrow-down-left", D: "bi bi-bag", I: "bi bi-graph-up-arrow" }[item.value.category?.type] || "bi bi-tag"));
const categoryLabel = computed(() => ({ R: "Receita", D: "Despesa", I: "Investimento" }[item.value.category?.type] || "Categoria"));
const valueClass = computed(() => item.value.category?.type === "R" ? "text-success" : item.value.category?.type === "D" ? "text-danger" : "text-primary");

function onCancel() {
  if (editingExisting.value) {
    item.value = props.budgetItem.clone();
    editingExisting.value = false;
    editMode.value = false;
  } else {
    emit("cancel-item", props.budgetItem);
  }
}

async function onDeleteClick() {
  if (!await dialog.show()) return;
  loading.show();
  try {
    await budgetService.delete(item.value.id);
    emit("delete-item", item.value);
    toast.success("Categoria excluída com sucesso!", { position: "top-center" });
  } finally {
    loading.hide();
  }
}

function onEditClick() {
  editingExisting.value = true;
  editMode.value = true;
}

async function onSaveClick() {
  if (!item.value.category?.id) {
    toast.warning("Selecione uma categoria.", { position: "top-center" });
    return;
  }

  loading.show();
  try {
    const creating = !item.value.id;
    const response = creating
      ? await budgetService.create(item.value)
      : await budgetService.modify(item.value.id, item.value);

    if (creating && response?.data) item.value = new item.value.constructor(response.data, props.options);
    editMode.value = false;
    editingExisting.value = false;
    emit("change-item", item.value.clone(), props.budgetItem);
    toast.success(`Categoria ${creating ? "incluída" : "atualizada"} com sucesso!`, { position: "top-center" });
  } finally {
    loading.hide();
  }
}
</script>

<style scoped>
td {
  padding: .9rem 1rem;
  border-color: var(--bs-border-color-translucent);
}

.editing-row {
  --bs-table-bg: rgba(var(--bs-primary-rgb), .035);
}

.category-icon {
  width: 2.25rem;
  height: 2.25rem;
  display: grid;
  place-items: center;
  flex: 0 0 auto;
  border-radius: .65rem;
  background: var(--bs-light);
  color: var(--bs-secondary);
}

.category-icon.is-income {
  color: var(--bs-success);
  background: rgba(var(--bs-success-rgb), .12);
}

.category-icon.is-expense {
  color: var(--bs-danger);
  background: rgba(var(--bs-danger-rgb), .1);
}

.category-icon.is-investment {
  color: var(--bs-primary);
  background: rgba(var(--bs-primary-rgb), .12);
}

.value-column {
  min-width: 12rem;
}

.value-input {
  max-width: 12rem;
}

.min-w-0 {
  min-width: 0;
}

@media (max-width: 575.98px) {
  td {
    padding: .75rem;
  }

  .category-icon {
    display: none;
  }
}
</style>