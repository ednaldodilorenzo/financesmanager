<template>
  <bootstrap-modal-screen
    :onClose="onCancelModal"
    :onConfirm="onSubmit"
    :title="`${form.id ? 'Alterar' : 'Nova'} Conta`"
  >
    <form
      :class="{
        'was-validated': v$.$dirty ? true : false,
      }"
      @submit.stop.prevent="onSubmit"
      id="frmAccount"
      class="row g-3 mb-3"
      autocomplete="off"
      novalidate
    >
      <div class="col-md-6">
        <bootstrap-input
          required-message="Por favor preencha o nome da conta"
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
            conta"
          :required="true"
          v-model="form.type"
          id="slcTipo"
          :options="[
            { id: 'A', description: 'Conta Corrente' },
            { id: 'C', description: 'Cartão de Crédito' },
            { id: 'D', description: 'Dinheiro' },
            { id: 'I', description: 'Investimento' },
          ]"
          :keyField="'id'"
          :valueField="'description'"
          label="Tipo"
        />
      </div>
      <div class="col-md-6">
        <bootstrap-select
          v-if="form.type === 'C'"
          required-message="Por favor preencha o dia de pagamento"
          :required="true"
          v-model="form.dueDay"
          id="slcDueDay"
          :options="dueDays"
          :keyField="'id'"
          :valueField="'description'"
          label="Dia do Pagamento"
        />
      </div>
    </form>
  </bootstrap-modal-screen>
</template>
<script setup>
import { ref, computed } from "vue";
import bootstrapModalScreen from "@/components/bootstrap-modal-screen.vue";
import BootstrapInput from "@/components/bootstrap-input.vue";
import BootstrapSelect from "@/components/bootstrap-select.vue";
import { useVuelidate } from "@vuelidate/core";
import accountService from "./account.service";
import { required, requiredIf } from "@vuelidate/validators";
import { useToast } from "vue-toastification";
import { useLoadingScreen } from "@/components/loading/useLoadingScreen";

const toast = useToast();
const loading = useLoadingScreen();
const form = ref({});
const dueDays = Array.from({ length: 28 }, (_, i) => ({
  id: 1 + i,
  description: i + 1 < 10 ? "0" + (i + 1) : "" + (i + 1),
}));

const rules = computed(() => ({
  name: {
    required,
  },
  type: {
    required,
  },
  dueDay: {
    required: requiredIf(() => {
      return form.value.type === "C";
    }),
  },
}));

const props = defineProps({
  onSaveModal: Function,
  onCancelModal: Function,
  item: {
    type: Object,
    required: false,
    default: () => ({
      name: "",
      type: "",
      dueDay: "",
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
  const payload = { name: form.value.name, type: form.value.type };
  if (form.value.dueDay) {
    payload.dueDay = form.value.dueDay;
  }
  const method = form.value.id
    ? accountService.modify(form.value.id, payload)
    : accountService.create(payload);

  method
    .then(() => {
      toast.success(
        `Conta ${form.value.id ? "atualizada" : "criada"} com sucesso!`,
        {
          position: "top-center",
        }
      );
      v$.value.$reset();
      props.onSaveModal();
    })
    .catch((e) => {
      toast.error("Falha na execução da solicitação!", {
        position: "top-center",
      });
    })
    .finally(() => {
      loading.hide();
    });
};
</script>
