<template>
  <bootstrap-modal-screen
    :onClose="onCancelModal"
    :onConfirm="onSubmit"
    :title="`${form.id ? 'Alterar' : 'Nova'} Transação`"
  >
    <form id="frmTransaction" class="row g-3 mb-3" autocomplete="off">
      <div class="col-md-6">
        <bootstrap-input
          type="date"
          required-message="Por favor preencha a data da transação"
          :required="true"
          v-model="form.transactionDate"
          id="iptDate"
          label="Data da Transação"
          name="date"
        />
      </div>
      <div class="col-md-6">
        <bootstrap-input
          type="date"
          required-message="Por favor preencha a data da transação"
          :required="true"
          v-model="form.paymentDate"
          id="iptPaymentDate"
          label="Data do Pagamento"
          name="date"
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
        <bootstrap-searcheable-select
          id="iptConta"
          label="Conta"
          v-model="form.account"
          display-field="name"
          value-field="id"
          :options="allAccounts"
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
import { formatDateUTC } from "@/utils/date";
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

// This is done because this component is loaded from another app.
const vCurrency = CurrencyDirective;

const props = defineProps({
  onSaveModal: Function,
  onCancelModal: Function,
  item: {
    type: Object,
    required: false,
    default: () => ({
      transactionDate: formatDateUTC(new Date(), "yyyy-MM-dd"),
      paymentDate: formatDateUTC(new Date(), "yyyy-MM-dd"),
      description: "",
      categoryId: undefined,
      category: null,
      accountId: undefined,
      account: null,
      value: undefined,
      detail: "",
      tags: [],
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

function getDependencies() {
  Promise.allSettled([
    categoryService.findAll({ paginate: false }),
    accountService.findAll(),
  ]).then((results) => {
    const [respCategories, respAccounts] = results;

    allCategories.value = respCategories.value.items;
    allAccounts.value = respAccounts.value.items;
    if (props.item.id) {
      transactionService.findById(props.item.id).then((resp) => {
        form.value = {
          ...resp.item,
          value: formatCurrency("" + resp.item.value),
          tags: resp.item.tags.map((value) => value.tag),
          transactionDate: formatDateUTC(
            resp.item.transactionDate,
            "yyyy-MM-dd"
          ),
          paymentDate: formatDateUTC(resp.item.paymentDate, "yyyy-MM-dd"),
        };
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
    .then((resp) => resp.items.map((item) => item.tag));
};

const onSubmit = () => {
  const { account, category, ...payload } = {
    ...form.value,
    value: parseCurrencyToNumber(form.value.value),
    paymentDate: new Date(form.value.paymentDate).toISOString(),
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
