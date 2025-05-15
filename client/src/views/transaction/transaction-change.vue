<template>
  <bootstrap-modal-screen
    :onClose="onCancelModal"
    :onConfirm="onSubmit"
    :title="`${form.id ? 'Alterar' : 'Nova'} Transação`"
  >
    <form id="frmTransaction" class="row g-3 mb-3" autocomplete="off">
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
        <bootstrap-input
          required-message="Por favor preencha a descrição da transação"
          :required="true"
          v-model="form.description"
          id="iptDescription"
          label="Descrição"
          name="description"
        />
      </div>
      <div
        :class="{
          'col-md-6': form.account?.type === 'C',
          'col-md-12': form.account?.type !== 'C',
        }"
      >
        <bootstrap-searcheable-select
          id="iptConta"
          label="Conta"
          v-model="form.account"
          display-field="name"
          value-field="id"
          :options="allAccounts"
        ></bootstrap-searcheable-select>
      </div>
      <div class="col-md-6" v-if="form.account?.type === 'C'">
        <bootstrap-select
          :options="creditCardDates"
          required-message="Por favor preencha a data da transação"
          :required="true"
          key-field="monthYear"
          value-field="description"
          v-model="form.monthYear"
          label="Fatura do Pagamento"
        />
      </div>
      <div class="col-md-6">
        <bootstrap-searcheable-select
          id="iptCategoria"
          label="Categoria"
          v-model="form.category"
          display-field="name"
          value-field="id"
          :options="allCategories"
        ></bootstrap-searcheable-select>
      </div>
      <div class="col-md-6">
        <bootstrap-input
          required-message="Por favor preencha o valor da transação"
          :required="true"
          v-model="form.value"
          v-currency
          id="iptValue"
          label="Valor"
          name="value"
        />
      </div>
      <div class="col-md-12">
        <bootstrap-select-tag
          v-model="form.tags"
          :options="searchTags"
          label="Tags"
          id="iptValue"
          name="value"
        />
      </div>
      <div class="col-md-12">
        <bootstrap-text-area
          :required="true"
          v-model="form.detail"
          label="Anotação"
          name="detail"
        />
      </div>
    </form>
  </bootstrap-modal-screen>
</template>
<script setup>
import { ref } from "vue";
import BootstrapInput from "@/components/bootstrap-input.vue";
import BootstrapSearcheableSelect from "@/components/bootstrap-searcheable-select.vue";
import BootstrapTextArea from "@/components/bootstrap-textarea.vue";
import BootstrapModalScreen from "@/components/bootstrap-modal-screen.vue";
import BootstrapSelectTag from "@/components/bootstrap-select-tag.vue";
import BootstrapSelect from "@/components/bootstrap-select.vue";
import { formatDateUTC, getAproximateMonths } from "@/utils/date";
import { required } from "@vuelidate/validators";
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

const props = defineProps({
  onSaveModal: Function,
  onCancelModal: Function,
  item: {
    type: Object,
    required: false,
    default: () => ({
      transactionDate: formatDateUTC(new Date(), "yyyy-MM-dd"),
      description: "",
      categoryId: undefined,
      category: null,
      accountId: undefined,
      account: null,
      value: undefined,
      detail: "",
      tags: [],
      monthYear: "",
    }),
  },
});

const form = ref({});

let allCategories = ref([]);
let allAccounts = ref([]);

const rules = {
  transactionDate: { required },
  paymentDate: { required },
  description: { required },
  category: { required },
  account: { required },
};

function onTransactionChange(event) {
  creditCardDates.value = getAproximateMonths(new Date()).map((item) => ({
    ...item,
    monthYear: `${item.month}/${item.year}`,
  }));
}

function getDependencies() {
  Promise.allSettled([
    categoryService.findAll({ paginate: false }),
    accountService.findAll({ paginate: false }),
  ]).then((results) => {
    const [respCategories, respAccounts] = results;

    allCategories.value = respCategories.value.data;
    allAccounts.value = respAccounts.value.data;
    if (props.item.id) {
      transactionService.findById(props.item.id).then((resp) => {
        let itemMonthYear = formatDateUTC(resp.data.paymentDate, "MM/yyyy");
        itemMonthYear =
          itemMonthYear[0] === "0" ? itemMonthYear.substring(1) : itemMonthYear;
        form.value = {
          ...resp.data,
          value: formatCurrency("" + resp.data.value),
          tags: resp.data.tags.map((value) => value.tag),
          transactionDate: formatDateUTC(
            resp.data.transactionDate,
            "yyyy-MM-dd"
          ),
          monthYear: itemMonthYear,
        };
        expense.value = resp.data.value < 0;
      });
    } else {
      form.value = props.item;
    }
  });
}

getDependencies();

const searchTags = (filter) => {
  return tagService
    .findAll({ filter: filter })
    .then((resp) => resp.data.map((item) => item.tag));
};

const onSubmit = () => {
  const { account, category, ...payload } = {
    ...form.value,
    value: expense.value
      ? -parseCurrencyToNumber(form.value.value)
      : parseCurrencyToNumber(form.value.value),
    paymentMonth: parseInt(
      form.value.monthYear.substring(0, form.value.monthYear.indexOf("/"))
    ),
    paymentYear: parseInt(
      form.value.monthYear.substring(
        form.value.monthYear.indexOf("/") + 1,
        form.value.monthYear.length
      )
    ),
    transactionDate: new Date(form.value.transactionDate).toISOString(),
    tags: form.value.tags.map((item) => ({
      tag: item,
    })),
    categoryId: form.value.category.id,
    accountId: form.value.account.id,
  };

  const method = form.value.id
    ? transactionService.modify(form.value.id, payload)
    : transactionService.create(payload);

  method.then(() => {
    toast.success(
      `Transação ${form.value.id ? "atualizada" : "criada"} com sucesso!`,
      {
        position: "top-center",
      }
    );
    props.onSaveModal();
  });
};
</script>
