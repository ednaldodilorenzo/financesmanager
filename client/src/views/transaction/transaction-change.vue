<template>
  <bootstrap-modal-screen
    :onClose="onCancelModal"
    :onConfirm="onSubmit"
    :title="`${form.id ? 'Alterar' : 'Nova'} Transação`"
  >
    <form
      id="frmTransaction"
      class="row g-3 mb-3"
      :class="{ 'was-validated': v$.$dirty }"
      autocomplete="off"
      novalidate
      @submit.prevent="onSubmit"
    >
      <div class="col-md-12">
        <div
          class="container d-flex justify-content-evenly align-items-center form-control"
        >
          <input
            type="radio"
            class="btn-check"
            name="options-outlined"
            id="success-outlined"
            autocomplete="off"
            :checked="!expense"
            @click="expense = false"
          />
          <label class="btn btn-outline-success" for="success-outlined"
            >Receita</label
          >
          <input
            type="radio"
            class="btn-check"
            name="options-outlined"
            id="danger-outlined"
            autocomplete="off"
            :checked="expense"
            @click="expense = true"
          />
          <label class="btn btn-outline-danger" for="danger-outlined"
            >Despesa</label
          >
        </div>
      </div>
      <div class="col-md-6">
        <bootstrap-input
          type="date"
          required-message="Por favor preencha a data da transação"
          :required="true"
          v-model="form.transactionDate"
          id="iptDate"
          label="Data da Transação"
          name="date"
          @change="onTransactionChange"
        />
      </div>
      <div class="col-md-6">
        <bootstrap-searcheable-select
          id="iptConta"
          label="Conta"
          v-model="form.account"
          display-field="name"
          value-field="id"
          :options="allAccounts"
          :required="true"
          required-message="Por favor selecione a conta"
        ></bootstrap-searcheable-select>
      </div>
      <div class="col-md-12" v-if="form.account?.type === 'C'">
        <bootstrap-select
          :options="creditCardDates"
          required-message="Por favor selecione a fatura do pagamento"
          :required="true"
          key-field="monthYear"
          value-field="description"
          v-model="form.monthYear"
          label="Fatura do Pagamento"
        />
      </div>
      <div class="col-md-12">
        <div class="d-flex flex-column gap-2">
          <div
            v-for="(split, index) in form.splits"
            :key="split.id || `new-${index}`"
            :class="{ 'split-card': form.splits.length > 1 }"
          >
            <div
              v-if="form.splits.length > 1"
              class="d-flex justify-content-between align-items-center px-2 mb-1"
            >
              <span class="small fw-semibold text-body-secondary"
                >Parte {{ index + 1 }}</span
              >
              <button
                type="button"
                class="btn btn-sm btn-link text-danger p-0"
                title="Remover parte"
                @click="removeSplit(index)"
              >
                <i class="bi bi-trash"></i>
              </button>
            </div>
            <div class="row g-2">
              <div class="col-md-5">
                <bootstrap-input
                  required-message="Por favor preencha a descrição"
                  :required="true"
                  v-model="split.description"
                  :id="`iptDescription${index}`"
                  label="Descrição"
                  :name="`description${index}`"
                />
              </div>
              <div class="col-md-4">
                <bootstrap-searcheable-select
                  :id="`iptCategoria${index}`"
                  label="Categoria"
                  v-model="split.category"
                  display-field="name"
                  value-field="id"
                  :options="allCategories"
                  :required="true"
                  required-message="Por favor selecione a categoria"
                ></bootstrap-searcheable-select>
              </div>
              <div class="col-md-3">
                <bootstrap-input
                  required-message="Por favor preencha o valor"
                  :required="true"
                  v-model="split.value"
                  v-currency
                  :id="`iptValue${index}`"
                  label="Valor"
                  :name="`value${index}`"
                />
              </div>
              <div class="col-md-6">
                <bootstrap-select-tag
                  v-model="split.tags"
                  :options="searchTags"
                  label="Tags"
                  :id="`iptTags${index}`"
                  :name="`tags${index}`"
                />
              </div>
              <div class="col-md-6">
                <bootstrap-text-area
                  v-model="split.detail"
                  label="Anotação"
                  rows="1"
                  :name="`detail${index}`"
                />
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="col-md-12">
        <div
          class="d-flex justify-content-between align-items-center border-top pt-3 mt-1"
        >
          <button
            type="button"
            class="btn btn-sm btn-outline-primary"
            @click="addSplit"
          >
            <i class="bi bi-plus-lg me-1"></i>Dividir transação
          </button>
          <span class="fw-semibold">
            Total: R$ {{ formatCurrency(String(totalValue)) }}
          </span>
        </div>
      </div>
    </form>
  </bootstrap-modal-screen>
</template>
<script setup>
import { computed, ref } from "vue";
import { useVuelidate } from "@vuelidate/core";
import { helpers, required, requiredIf } from "@vuelidate/validators";
import BootstrapInput from "@/components/bootstrap-input.vue";
import BootstrapSearcheableSelect from "@/components/bootstrap-searcheable-select.vue";
import BootstrapTextArea from "@/components/bootstrap-textarea.vue";
import BootstrapModalScreen from "@/components/bootstrap-modal-screen.vue";
import BootstrapSelectTag from "@/components/bootstrap-select-tag.vue";
import BootstrapSelect from "@/components/bootstrap-select.vue";
import { formatDateUTC, getAproximateMonths } from "@/utils/date";
import categoryService from "@/views/category/category.service";
import accountService from "@/views/account/account.service";
import tagService from "@/views/tag/tag.service";
import transactionService from "./transaction.service";
import { parseCurrencyToNumber } from "@/utils/numbers";
import CurrencyDirective from "@/components/currency.directive";
import { useToast } from "vue-toastification";
import { formatCurrency } from "@/utils/numbers";

const toast = useToast();
let creditCardDates = ref(
  getAproximateMonths(new Date()).map((item) => ({
    ...item,
    monthYear: `${item.month}/${item.year}`,
  }))
);

// This is done because this component is loaded from another app.
const vCurrency = CurrencyDirective;

const expense = ref(true);

const newSplit = () => ({
  id: undefined,
  description: "",
  category: null,
  value: "",
  detail: "",
  tags: [],
});

const props = defineProps({
  onSaveModal: Function,
  onCancelModal: Function,
  item: {
    type: Object,
    required: false,
    default: () => ({
      transactionDate: formatDateUTC(new Date(), "yyyy-MM-dd"),
      account: null,
      monthYear: "",
    }),
  },
});

const form = ref({ splits: [] });

let allCategories = ref([]);
let allAccounts = ref([]);

const totalValue = computed(() =>
  (form.value.splits || []).reduce(
    (total, split) => total + parseCurrencyToNumber(split.value),
    0
  )
);

const rules = computed(() => ({
  transactionDate: { required },
  account: { required },
  monthYear: { required: requiredIf(() => form.value.account?.type === "C") },
  splits: {
    required,
    $each: helpers.forEach({
      description: { required },
      category: { required },
      value: {
        required,
        nonZero: (value) => parseCurrencyToNumber(value) !== 0,
      },
    }),
  },
}));
const v$ = useVuelidate(rules, form);

function addSplit() {
  form.value.splits.push(newSplit());
}

function removeSplit(index) {
  form.value.splits.splice(index, 1);
}

function onTransactionChange(event) {
  creditCardDates.value = getAproximateMonths(new Date()).map((item) => ({
    ...item,
    monthYear: `${item.month}/${item.year}`,
  }));
}

function toMonthYear(date) {
  const monthYear = formatDateUTC(date, "MM/yyyy");
  return monthYear[0] === "0" ? monthYear.substring(1) : monthYear;
}

function getDependencies() {
  Promise.allSettled([
    categoryService.findAll({ paginate: false }),
    accountService.findAll({ paginate: false }),
  ]).then((results) => {
    const [respCategories, respAccounts] = results;

    allCategories.value = respCategories.value;
    allAccounts.value = respAccounts.value;
    if (props.item.transactionId) {
      transactionService
        .findWithSplits(props.item.transactionId)
        .then((resp) => {
          form.value = {
            id: resp.id,
            account: resp.account,
            transactionDate: formatDateUTC(resp.transactionDate, "yyyy-MM-dd"),
            monthYear: toMonthYear(resp.paymentDate),
            splits: resp.splits.map((split) => ({
              ...split,
              value: formatCurrency("" + split.value),
              tags: split.tags.map((value) => value.tag),
            })),
          };
          expense.value = resp.value < 0;
        });
    } else {
      form.value = { ...props.item, splits: [newSplit()] };
    }
  });
}

getDependencies();

const searchTags = (filter) => {
  return tagService
    .findAll({ filter: filter })
    .then((resp) => resp.map((item) => item.tag));
};

const signed = (value) => (expense.value ? -value : value);

const onSubmit = async () => {
  if (!(await v$.value.$validate())) return;

  const { monthYear, splits } = form.value;
  const header = {
    accountId: form.value.account.id,
    transactionDate: new Date(form.value.transactionDate).toISOString(),
    paymentMonth: parseInt(monthYear.substring(0, monthYear.indexOf("/"))),
    paymentYear: parseInt(monthYear.substring(monthYear.indexOf("/") + 1)),
  };
  const toSplitPayload = (split) => ({
    id: split.id,
    description: split.description,
    categoryId: split.category.id,
    value: signed(parseCurrencyToNumber(split.value)),
    detail: split.detail,
    tags: split.tags.map((item) => ({ tag: item })),
  });
  const payload = {
    ...header,
    id: form.value.id,
    value: signed(totalValue.value),
    splits: splits.map(toSplitPayload),
  };

  const method = form.value.id
    ? transactionService.modifySplit(form.value.id, payload)
    : transactionService.createSplit(payload);

  method.then(() => {
    toast.success(
      `Transação ${form.value.id ? "atualizada" : "criada"} com sucesso!`,
      {
        position: "top-center",
      }
    );
    v$.value.$reset();
    props.onSaveModal();
  });
};
</script>

<style scoped>
.split-card {
  border: 1px solid var(--bs-border-color);
  border-radius: 0.5rem;
  padding: 0.5rem 0;
}
</style>
