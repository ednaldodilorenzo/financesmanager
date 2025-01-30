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
  <div
    v-if="importedTransactions.length > 0"
    class="d-flex justify-content-start"
  >
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
  <div
    v-for="item in importedTransactions"
    :key="item.id"
    class="card my-3"
    :class="{ 'border-danger': item.duplicated }"
  >
    <div v-if="item.duplicated" class="card-header text-danger">Duplicado</div>
    <div class="card-body">
      <form class="row gy-2 gx-3 align-items-center">
        <div class="col-auto">
          <input
            class="form-check-input"
            type="checkbox"
            v-model="item.checked"
            id="flexCheckDefault"
          />
        </div>
        <div class="col-auto">
          <div class="form-floating">
            <input
              type="date"
              class="form-control"
              :value="item.formatted_date"
              disabled
              placeholder="Jane Doe"
            />
            <label>Data</label>
          </div>
        </div>
        <div class="col-auto">
          <div class="form-floating">
            <input
              type="text"
              :value="item.description"
              class="form-control"
              placeholder="Descrição"
            />
            <label for="floatingInput">Descrição</label>
          </div>
        </div>
        <div class="col-auto">
          <div class="form-floating">
            <select v-model="item.accountId" class="form-select">
              <option></option>
              <option
                :value="account.id"
                v-for="account in accounts"
                :key="account.id"
              >
                {{ account.name }}
              </option>
            </select>
            <label>Conta</label>
          </div>
        </div>
        <div class="col-auto">
          <div class="form-floating">
            <select v-model="item.categoryId" class="form-select">
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
        <div class="col-auto">
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
        <div class="col-auto">
          <div class="input-group">
            <span class="input-group-text">R$</span>
            <div class="form-floating">
              <input
                type="text"
                v-currency
                v-model="item.formatted_value"
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
              class="btn btn-primary"
              @click.prevent="importTransaction(item)"
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
<script setup>
import transactionService from "./transaction.service";
import categoryService from "../category/category.service";
import accountService from "../account/account.service";
import LoadingScreen from "@/components/loading-screen.vue";
import { format } from "date-fns";
import { ref } from "vue";
import { useToast } from "vue-toastification";
import { parseCurrencyToNumber } from "@/utils/numbers";

const toast = useToast();

let formData = ref({
  file: null, // File will be stored here
});

let fileType = ref("");

let categories = [],
  accounts = [];
let importedTransactions = ref([]);

const loading = ref(false);

function getDependencies() {
  loading.value = true;
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
      loading.value = false;
    });
}

getDependencies();

function sendBatchData(data) {
  const payload = data.map(({ category, account, clazz, ...item }) => ({
    ...item,
    value: parseCurrencyToNumber(item.formatted_value),
    paymentDate: new Date(item.paymentDate).toISOString(),
    transactionDate: new Date(item.transactionDate).toISOString(),
  }));
  loading.value = true;
  transactionService
    .sendBatchImport(payload)
    .then(() => {
      toast.success(`Transações criadas com sucesso!`, {
        position: "top-center",
      });
      importedTransactions.value = importedTransactions.value.filter(
        (value) => !data.includes(value)
      );
    })
    .finally(() => {
      loading.value = false;
    });
}

function importSelectedTransactions() {
  const filteredTransactions = importedTransactions.value.filter(
    (item) => item.checked
  );

  sendBatchData(filteredTransactions);
}

function importNonDuplicated() {
  const filteredTransactions = importedTransactions.value.filter(
    (item) => !item.duplicated
  );
  sendBatchData(filteredTransactions);
}

function discardDuplicated() {
  importedTransactions.value = importedTransactions.value.filter(
    (value) => !value.duplicated
  );
}

function importTransaction(itemClicked) {
  loading.value = true;

  const payload = {
    categoryId: itemClicked.categoryId,
    accountId: itemClicked.accountId,
    detail: itemClicked.detail,
    description: itemClicked.description,
    value: parseCurrencyToNumber(itemClicked.formatted_value),
    paymentDate: new Date(itemClicked.paymentDate).toISOString(),
    transactionDate: new Date(itemClicked.transactionDate).toISOString(),
  };

  transactionService
    .create(payload)
    .then(() => {
      toast.success(`Transação criada com sucesso!`, {
        position: "top-center",
        timeout: 5000,
        closeOnClick: true,
        pauseOnFocusLoss: true,
        pauseOnHover: true,
        draggable: true,
        draggablePercent: 0.6,
        showCloseButtonOnHover: false,
        hideProgressBar: true,
        closeButton: "button",
        icon: true,
        rtl: false,
      });
      importedTransactions.value = importedTransactions.value.filter(
        (value) => value != itemClicked
      );
    })
    .finally(() => {
      loading.value = false;
    });
}

function handleFileUpload(event) {
  formData.value.file = event.target.files[0]; // Capture the uploaded file
}

function removeItem(itemClicked) {
  importedTransactions.value = importedTransactions.value.filter(
    (value) => value !== itemClicked
  );
}

function submitForm(e) {
  loading.value = true;
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
      importedTransactions.value = resp.items.map((item) => ({
        ...item,
        formatted_date: format(item.paymentDate, "yyyy-MM-dd"),
        formatted_value: (item.value / 100).toLocaleString("pt-BR", {
          minimumFractionDigits: 2,
          maximumFractionDigits: 2,
        }),
        category: item.categoryId,
        account: item.accountId,
        checked: false,
        clazz: item.duplicated ? "table-danger" : null,
      }));
    })
    .finally(() => {
      loading.value = false;
    });
}
</script>
<style scoped>
.value-summary {
  border: solid 1px black;
  border-radius: 3.125rem;
}

.value-summary.expense {
  color: red;
  border-color: red;
}

.value-summary.earn {
  color: green;
  border-color: green;
}
</style>
