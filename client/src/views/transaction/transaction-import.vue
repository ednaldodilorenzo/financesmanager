<template>
  <div class="mt-3">
    <loading-screen :loading="loading" />
    <div class="d-flex justify-content-between">
      <h2 class="fs-4">Transações</h2>
      <nav style="--bs-breadcrumb-divider: '>'" aria-label="breadcrumb">
        <ol class="breadcrumb">
          <li class="breadcrumb-item"><a href="#">Home</a></li>
          <li class="breadcrumb-item active"><a href="#">Importar</a></li>
        </ol>
      </nav>
    </div>
  </div>
  <hr />
  <form class="d-flex justify-content-center my-3" @submit.prevent="submitForm">
    <select name="fileType" v-model="fileType" class="form-control mx-3">
      <option value="BBCA">Conta Corrente BB</option>
      <option value="C6CC">Cartão de Crédito C6</option>
      <option value="CUAL">Arquivo Customizado</option>
    </select>
    <input
      @change="handleFileUpload"
      name="file"
      class="form-control"
      type="file"
      id="formFile"
    />

    <select
      name="accountId"
      v-if="['BBCA', 'C6CC'].includes(fileType)"
      class="form-control mx-3"
      v-model="currentAccount"
    >
      <option></option>
      <option :value="account.id" v-for="account in accounts" :key="account.id">
        {{ account.name }}
      </option>
    </select>
    <input
      type="date"
      v-if="fileType === 'C6CC'"
      name="paymentDate"
      class="form-control mx-3"
    />
    <button type="submit" class="btn btn-primary mx-3">Enviar</button>
  </form>
  <div v-if="state.collection.length > 0" class="d-flex justify-content-center">
    <div class="col-auto">
      <button
        @click.prevent="importSelectedTransactions"
        class="btn btn-success mx-3"
      >
        Importar Selecionados
      </button>
    </div>
    <div class="col-auto">
      <button @click.prevent="importNonDuplicated" class="btn btn-primary mx-3">
        Importar Não Duplicados
      </button>
    </div>
    <div class="col-auto">
      <button @click.prevent="discardDuplicated" class="btn btn-danger mx-3">
        Descartar Duplicados
      </button>
    </div>
  </div>
  <ValidateEach
    v-for="item in state.collection"
    :key="item"
    :state="item"
    :rules="rules"
  >
    <template #default="{ v }">
      <div class="card my-3" :class="{ 'border-danger': item.duplicated }">
        <div v-if="item.duplicated" class="card-header text-danger">
          Duplicado
        </div>
        <div class="card-body">
          <form class="row gy-2 gx-3 align-items-center">
            <div class="col-auto d-flex align-items-center">
              <input
                class="form-check-input me-3"
                type="checkbox"
                v-model="item.checked"
                id="flexCheckDefault"
              />
              <div class="form-floating">
                <input
                  type="date"
                  class="form-control"
                  :value="v.formatted_date.$model"
                  disabled
                  placeholder="Jane Doe"
                />
                <label>Data</label>
              </div>
            </div>
            <div class="col">
              <div class="form-floating">
                <input
                  type="text"
                  v-model="v.description.$model"
                  class="form-control"
                  :class="{ 'is-invalid': v.description.$errors.length > 0 }"
                  placeholder="Descrição"
                />
                <label for="floatingInput">Descrição</label>
              </div>
            </div>
            <div class="col">
              <div class="form-floating">
                <select
                  v-model="v.category.$model"
                  class="form-select"
                  :class="{ 'is-invalid': v.category.$errors.length > 0 }"
                >
                  <option></option>
                  <option
                    :value="category.id"
                    v-for="category in categories"
                    :key="category.id"
                  >
                    {{ category.name }}
                  </option>
                </select>
                <label>Categoria</label>
              </div>
            </div>
            <div class="col">
              <div class="form-floating">
                <input
                  type="text"
                  v-model="item.detail"
                  class="form-control"
                  placeholder="Anotação"
                />
                <label for="floatingInput">Anotação</label>
              </div>
            </div>
            <div class="col">
              <div class="form-floating">
                <bootstrap-select-tag
                  v-model="item.tags"
                  :options="searchTags"
                  :floating="true"
                  id="iptTags"
                  name="tags"
                />
              </div>
            </div>
            <div class="col-2">
              <div class="input-group">
                <span class="input-group-text">R$</span>
                <div class="form-floating">
                  <input
                    type="text"
                    v-currency
                    v-model="v.formatted_value.$model"
                    :class="{
                      'is-invalid': v.formatted_value.$errors.length > 0,
                    }"
                    class="form-control"
                    placeholder="Valor"
                  />
                  <label>Valor</label>
                </div>
              </div>
            </div>
            <div class="col-auto">
              <div
                class="btn-group"
                role="group"
                aria-label="Basic outlined example"
              >
                <button
                  type="button"
                  class="btn btn-primary"
                  @click.prevent="importTransaction(v, item)"
                >
                  <i class="bi bi-arrow-down-circle-fill"></i>
                </button>
              </div>
              <div
                class="btn-group"
                role="group"
                aria-label="Basic outlined example"
              >
                <button
                  type="button"
                  class="btn btn-danger mx-3"
                  @click.prevent="removeItem(item)"
                >
                  <i class="bi bi-x-circle-fill"></i>
                </button>
              </div>
            </div>
          </form>
        </div>
      </div>
    </template>
  </ValidateEach>
</template>

<script setup>
import { ref } from "vue";
import { useVuelidate } from "@vuelidate/core";
import { required } from "@vuelidate/validators";
import { ValidateEach } from "@vuelidate/components";
import transactionService from "./transaction.service";
import { useLoadingScreen } from "@/components/loading/useLoadingScreen";
import { format } from "date-fns";
import categoryService from "../category/category.service";
import accountService from "../account/account.service";
import { parseCurrencyToNumber } from "@/utils/numbers";
import { useToast } from "vue-toastification";
import { currencyBRL } from "@/components/filters/currency.filter";
import BootstrapSelectTag from "@/components/bootstrap-select-tag.vue";
import tagService from "@/views/tag/tag.service";

const loading = useLoadingScreen();
const toast = useToast();
let currentAccount = ref(null);

let categories = [],
  accounts = [];

const state = ref({
  collection: [],
});

let fileType = ref("");

function getDependencies() {
  loading.show();
  Promise.allSettled([
    categoryService.findAll({ paginate: false }),
    accountService.findAll(),
  ])
    .then((results) => {
      const [respCategories, respAccounts] = results;
      categories = respCategories.value.items;
      accounts = respAccounts.value.items;
    })
    .finally(() => {
      loading.hide();
    });
}

getDependencies();

const searchTags = (filter) => {
  return tagService
    .findAll({ filter: filter })
    .then((resp) => resp.items.map((item) => item.tag));
};

const rules = {
  formatted_date: {
    required,
  },
  description: {
    required,
  },
  category: {
    required,
  },
  formatted_value: {
    required,
  },
};

function removeItem(itemClicked) {
  state.value.collection = state.value.collection.filter(
    (value) => value !== itemClicked
  );
}

let v$ = useVuelidate();

function validate() {
  return v$.value.$validate();
}

function importNonDuplicated() {
  validate().then((result) => {
    if (result) {
      const filteredTransactions = state.value.collection.filter(
        (item) => !item.duplicated
      );
      sendBatchData(filteredTransactions);
    }
  });
}

function importTransaction(validator, itemClicked) {
  validator.$validate().then((isValid) => {
    if (isValid) {
      loading.show();

      const payload = {
        categoryId: itemClicked.category,
        accountId: currentAccount.value,
        detail: itemClicked.detail,
        description: itemClicked.description,
        tags: itemClicked.tags.map((item) => ({
          tag: item,
        })),
        value: parseCurrencyToNumber(itemClicked.formatted_value),
        paymentDate: new Date(itemClicked.paymentDate).toISOString(),
        transactionDate: new Date(itemClicked.transactionDate).toISOString(),
      };

      transactionService
        .create(payload)
        .then(() => {
          toast.success(`Transação criada com sucesso!`, {
            position: "top-center",
          });
          state.value.collection = state.value.collection.filter(
            (value) => value != itemClicked
          );
        })
        .finally(() => {
          loading.hide();
        });
    }
  });
}

function discardDuplicated() {
  state.value.collection = state.value.collection.filter(
    (value) => !value.duplicated
  );
}

function sendBatchData(data) {
  const payload = data.map(({ category, account, clazz, ...item }) => ({
    ...item,
    value: parseCurrencyToNumber(item.formatted_value),
    categoryId: category,
    accountId: currentAccount.value,
    tags: item.tags.map((value) => ({
      tag: value,
    })),
    paymentDate: new Date(item.paymentDate).toISOString(),
    transactionDate: new Date(item.transactionDate).toISOString(),
  }));
  loading.show();
  transactionService
    .sendBatchImport(payload)
    .then(() => {
      toast.success(`Transações criadas com sucesso!`, {
        position: "top-center",
      });
      state.value.collection = state.value.collection.filter(
        (value) => !data.includes(value)
      );
    })
    .finally(() => {
      loading.hide();
    });
}

let formData = ref({
  file: null, // File will be stored here
});

function handleFileUpload(event) {
  formData.value.file = event.target.files[0]; // Capture the uploaded file
}

function submitForm(e) {
  loading.show();
  const payload = new FormData();
  payload.append("file", formData.value.file); // Add the file
  if (e.target.elements.accountId) {
    payload.append("accountId", e.target.elements.accountId.value);
  }
  if (e.target.elements.paymentDate) {
    payload.append("paymentDate", e.target.elements.paymentDate.value);
  }
  payload.append("fileType", e.target.elements.fileType.value);
  transactionService
    .prepareForImport(payload)
    .then((resp) => {
      state.value.collection = resp.items.map((item) => ({
        ...item,
        formatted_date: format(item.paymentDate, "yyyy-MM-dd"),
        formatted_value: currencyBRL(item.value),
        category: item.categoryId ? item.categoryId : "",
        checked: false,
        tags: [],
        clazz: item.duplicated ? "table-danger" : null,
      }));
    })
    .finally(() => {
      loading.hide();
    });
}
</script>

<style>
.error {
  color: tomato;
}
</style>
