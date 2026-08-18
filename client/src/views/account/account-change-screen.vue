<template>
  <BootstrapModalScreen :onClose="onCancelModal" :onConfirm="onSubmit" :title="`${form.id ? 'Alterar' : 'Nova'} conta`">
    <form id="frmAccount" class="mb-3" autocomplete="off" novalidate @submit.prevent="onSubmit">
      <div class="modal-intro d-flex align-items-center gap-3 p-3 mb-4">
        <span class="intro-icon"><i :class="['bi', selectedType.icon]"></i></span>
        <div>
          <strong class="d-block">{{ form.id ? 'Editar dados da conta' : 'Cadastrar uma nova conta' }}</strong>
          <span class="small text-body-secondary">Use um nome fácil de reconhecer nos lançamentos.</span>
        </div>
      </div>

      <div class="row g-3">
        <div class="col-12">
          <BootstrapInput v-model="form.name" id="iptNome" name="nome" label="Nome da conta" :required="true"
            required-message="Por favor, preencha o nome da conta" />
        </div>
        <div class="col-md-7">
          <BootstrapSelect v-model="form.type" id="slcTipo" label="Tipo" :required="true"
            required-message="Por favor, selecione o tipo da conta" :options="accountTypes" keyField="id"
            valueField="description" />
        </div>
        <div v-if="form.type === 'C'" class="col-md-5">
          <BootstrapSelect v-model="form.dueDay" id="slcDueDay" label="Dia do pagamento" :required="true"
            required-message="Informe o dia do pagamento" :options="dueDays" keyField="id" valueField="description" />
        </div>
      </div>
    </form>
  </BootstrapModalScreen>
</template>

<script setup>
import { computed, ref, watch } from "vue";
import { useVuelidate } from "@vuelidate/core";
import { required, requiredIf } from "@vuelidate/validators";
import { useToast } from "vue-toastification";
import BootstrapModalScreen from "@/components/bootstrap-modal-screen.vue";
import BootstrapInput from "@/components/bootstrap-input.vue";
import BootstrapSelect from "@/components/bootstrap-select.vue";
import { useLoadingScreen } from "@/components/loading/useLoadingScreen";
import accountService from "./account.service";

const props = defineProps({
  onSaveModal: { type: Function, required: true },
  onCancelModal: { type: Function, required: true },
  item: { type: Object, default: () => ({ name: "", type: "", dueDay: "" }) },
});

const toast = useToast();
const loading = useLoadingScreen();
const form = ref({ ...props.item });
const accountTypes = [
  { id: "A", description: "Conta Corrente", icon: "bi-bank" },
  { id: "C", description: "Cartão de Crédito", icon: "bi-credit-card" },
  { id: "D", description: "Dinheiro", icon: "bi-cash-stack" },
  { id: "I", description: "Investimento", icon: "bi-graph-up-arrow" },
];
const dueDays = Array.from({ length: 28 }, (_, index) => ({ id: index + 1, description: String(index + 1).padStart(2, "0") }));
const selectedType = computed(() => accountTypes.find(({ id }) => id === form.value.type) || { icon: "bi-wallet2" });
const rules = computed(() => ({
  name: { required },
  type: { required },
  dueDay: { required: requiredIf(() => form.value.type === "C") },
}));
const v$ = useVuelidate(rules, form);

watch(() => form.value.type, (type) => {
  if (type !== "C") form.value.dueDay = "";
});

async function onSubmit() {
  if (!await v$.value.$validate()) return;
  loading.show();
  try {
    const payload = { name: form.value.name.trim(), type: form.value.type };
    if (form.value.type === "C") payload.dueDay = form.value.dueDay;
    if (form.value.id) await accountService.modify(form.value.id, payload);
    else await accountService.create(payload);
    toast.success(`Conta ${form.value.id ? "atualizada" : "criada"} com sucesso!`, { position: "top-center" });
    v$.value.$reset();
    props.onSaveModal();
  } catch (error) {
    toast.error("Falha na execução da solicitação!", { position: "top-center" });
  } finally { loading.hide(); }
}
</script>

<style scoped>
.modal-intro {
  border-radius: .75rem;
  background: var(--bs-tertiary-bg);
}

.intro-icon {
  width: 2.75rem;
  height: 2.75rem;
  display: grid;
  place-items: center;
  flex: 0 0 auto;
  border-radius: .8rem;
  background: rgba(var(--bs-primary-rgb), .1);
  color: var(--bs-primary);
  font-size: 1.15rem;
}
</style>