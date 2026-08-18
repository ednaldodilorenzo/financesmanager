<template>
  <BootstrapModalScreen :onClose="onCancelModal" :onConfirm="onSubmit"
    :title="`${form.id ? 'Alterar' : 'Nova'} categoria`">
    <form id="frmCategory" class="mb-3" autocomplete="off" novalidate @submit.prevent="onSubmit">
      <div class="modal-intro d-flex align-items-center gap-3 p-3 mb-4">
        <span :class="['intro-icon', selectedType.tone]"><i :class="['bi', selectedType.icon]"></i></span>
        <div><strong class="d-block">{{ form.id ? 'Editar categoria' : 'Cadastrar uma categoria' }}</strong><span
            class="small text-body-secondary">Categorias ajudam a organizar seus relatórios e planejamentos.</span>
        </div>
      </div>
      <div class="row g-3">
        <div class="col-md-7">
          <BootstrapInput v-model="form.name" id="iptNome" name="nome" label="Nome da categoria" :required="true"
            required-message="Por favor, preencha o nome da categoria" />
        </div>
        <div class="col-md-5">
          <BootstrapSelect v-model="form.type" id="slcTipo" label="Tipo" :required="true"
            required-message="Por favor, selecione o tipo da categoria" :options="categoryTypes" keyField="id"
            valueField="description" />
        </div>
      </div>
    </form>
  </BootstrapModalScreen>
</template>

<script setup>
import { computed, ref } from "vue";
import { useVuelidate } from "@vuelidate/core";
import { required } from "@vuelidate/validators";
import { useToast } from "vue-toastification";
import BootstrapModalScreen from "@/components/bootstrap-modal-screen.vue";
import BootstrapInput from "@/components/bootstrap-input.vue";
import BootstrapSelect from "@/components/bootstrap-select.vue";
import { useLoadingScreen } from "@/components/loading/useLoadingScreen";
import categoryService from "./category.service";

const props = defineProps({
  onSaveModal: { type: Function, required: true },
  onCancelModal: { type: Function, required: true },
  item: { type: Object, default: () => ({ name: "", type: "" }) },
});
const toast = useToast();
const loading = useLoadingScreen();
const form = ref({ ...props.item });
const categoryTypes = [
  { id: "R", description: "Receita", icon: "bi-arrow-down-left", tone: "is-income" },
  { id: "D", description: "Despesa", icon: "bi-arrow-up-right", tone: "is-expense" },
  { id: "I", description: "Investimento", icon: "bi-graph-up-arrow", tone: "is-investment" },
];
const selectedType = computed(() => categoryTypes.find(({ id }) => id === form.value.type) || { icon: "bi-tag", tone: "" });
const v$ = useVuelidate({ name: { required }, type: { required } }, form);

async function onSubmit() {
  if (!await v$.value.$validate()) return;
  loading.show();
  try {
    const payload = { name: form.value.name.trim(), type: form.value.type };
    if (form.value.id) await categoryService.modify(form.value.id, payload);
    else await categoryService.create(payload);
    toast.success(`Categoria ${form.value.id ? "atualizada" : "criada"} com sucesso!`, { position: "top-center" });
    v$.value.$reset();
    props.onSaveModal();
  } catch (error) {
    toast.error("Não foi possível salvar a categoria.", { position: "top-center" });
  } finally { loading.hide(); }
}
</script>

<style scoped>
.modal-intro {
  border-radius: .75rem;
  background: var(--bs-tertiary-bg)
}

.intro-icon {
  width: 2.75rem;
  height: 2.75rem;
  display: grid;
  place-items: center;
  flex: 0 0 auto;
  border-radius: .8rem;
  font-size: 1.15rem
}

.intro-icon:not(.is-income):not(.is-expense):not(.is-investment) {
  background: rgba(var(--bs-primary-rgb), .1);
  color: var(--bs-primary)
}
</style>