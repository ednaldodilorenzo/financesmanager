<template>
  <bootstrap-modal-screen
    :onClose="onCancelModal"
    :onConfirm="onSubmit"
    :title="`${form.id ? 'Alterar' : 'Nova'} Categoria`"
  >
    <form
      :class="{
        'was-validated': v$.$dirty ? true : false,
      }"
      @submit.stop.prevent="onSubmit"
      id="frmCategory"
      class="row 3 mb-3"
      autocomplete="off"
      novalidate
    >
      <div class="col-md-6">
        <bootstrap-input
          required-message="Por favor preencha o nome da categoria"
          :required="true"
          v-model="form.name"
          id="iptNome"
          label="Nome"
          name="nome"
        />
      </div>
      <div class="col-md-6">
        <bootstrap-select
          required-message="Por favor preencha o tipo da
          categoria"
          :required="true"
          v-model="form.type"
          id="slcTipo"
          :options="[
            { id: '', description: 'Selecione...' },
            { id: 'R', description: 'Receita' },
            { id: 'D', description: 'Despesa' },
            { id: 'I', description: 'Investimento' },
          ]"
          :keyField="'id'"
          :valueField="'description'"
          label="Tipo"
        />
      </div>
    </form>
  </bootstrap-modal-screen>
</template>
<script setup>
import { ref } from "vue";
import BootstrapModalScreen from "@/components/bootstrap-modal-screen.vue";
import BootstrapInput from "@/components/bootstrap-input.vue";
import BootstrapSelect from "@/components/bootstrap-select.vue";
import { useVuelidate } from "@vuelidate/core";
import categoryService from "./category.service";
import { required } from "@vuelidate/validators";
import { useToast } from "vue-toastification";
import { useLoadingScreen } from "@/components/loading/useLoadingScreen";

const toast = useToast();
const loading = useLoadingScreen();
const form = ref({});

const rules = {
  name: {
    required,
  },
  type: {
    required,
  },
};

const props = defineProps({
  onSaveModal: Function,
  onCancelModal: Function,
  item: {
    type: Object,
    required: false,
    default: () => ({
      name: "",
      type: "",
    }),
  },
});

form.value = props.item;

const v$ = useVuelidate(rules, form.value);

const onSubmit = () => {
  v$.value.$validate();

  if (v$.value.$error) {
    return;
  }

  loading.show();
  const payload = { ...form.value, type: form.value.type };
  const method = form.value.id
    ? categoryService.modify(form.value.id, payload)
    : categoryService.create(payload);

  method
    .then(() => {
      toast.success(
        `Categoria ${form.value.id ? "atualizada" : "criada"} com sucesso!`,
        {
          position: "top-center",
        }
      );
      v$.value.$reset();
      props.onSaveModal();
    })
    .catch((e) => {
      console.error(e);
    })
    .finally(() => {
      loading.hide();
    });
};
</script>
