<template>
  <div class="mt-3">
    <div class="d-flex justify-content-between">
      <h2 class="fs-4">Transações</h2>
      <nav style="--bs-breadcrumb-divider: '>'" aria-label="breadcrumb">
        <ol class="breadcrumb">
          <li class="breadcrumb-item"><a href="#">Home</a></li>
          <li class="breadcrumb-item active"><a href="#">Transações</a></li>
        </ol>
      </nav>
    </div>
  </div>
  <div class="card mb-3">
    <div class="card-body p-2">
      <div class="d-flex justify-content-center my-3">
        <Calendar @date-change="onChangeDebounced"></Calendar>
      </div>
    </div>
  </div>
  <summary-data
    :data="[
      {
        title: 'Receitas',
        value: summary.earns,
      },
      {
        title: 'Despesas',
        value: Math.abs(summary.expenses),
        percent: Math.abs(summary.expenses / summary.earns),
        percentMessage: 'do Total de Receitas',
      },
      {
        title: 'Saldo',
        value: summary.earns + summary.expenses,
      },
      {
        title: 'Investido',
        value: Math.abs(summary.investments),
      },
    ]"
  />
  <div class="card">
    <div class="card-body p-2">
      <nav class="navbar bg-body-tertiary">
        <div class="d-inline-flex w-100 justify-content-around">
          <div class="d-flex align-items-center" style="width: 30%">
            <label for="exampleDataList" class="form-label me-3">Tipo</label>
            <select @change="onTypeChange" class="form-control">
              <option value=""></option>
              <option value="R">Receita</option>
              <option value="D">Despesa</option>
            </select>
          </div>
          <div class="d-flex align-items-center" style="width: 30%">
            <label for="exampleDataList" class="form-label me-3">Conta</label>
            <select @change="onAccountChange" class="form-control">
              <option></option>
              <option v-for="account in accounts" :key="account">
                {{ account.name }}
              </option>
            </select>
          </div>
          <div class="d-flex align-items-center" style="width: 30%">
            <label for="exampleDataList" class="form-label me-3"
              >Categoria</label
            >
            <bootstrap-searcheable-select
              displayField="name"
              valueField="id"
              v-model="selectedCategory"
              :options="categories"
            ></bootstrap-searcheable-select>
          </div>
        </div>
      </nav>
      <Table
        :fields="fields"
        :showCSVButton="true"
        :csvHandler="exportToCSV"
        :showPagination="false"
        :items="filteredItems"
        :showFilter="true"
        :actions="[
          {
            name: 'edit',
            title: 'Editar Transação',
            icon: 'bi bi-pencil-fill',
            clazz: 'link-primary',
            handler: handleEdit,
          },
          {
            name: 'delete',
            title: 'Excluir Transação',
            icon: 'bi bi-trash-fill',
            clazz: 'link-danger',
            handler: handleDelete,
          },
        ]"
        @trigger-page="getList"
        @new-clicked="onNewClicked()"
      ></Table>
    </div>
  </div>
</template>
<script setup>
import Table from "@/components/bootstrap-table.vue";
import BootstrapSearcheableSelect from "@/components/bootstrap-searcheable-select.vue";
import Calendar from "@/components/bootstrap-calendar.vue";
import transactionService from "./transaction.service";
import categoryService from "../category/category.service";
import accountService from "../account/account.service";
import { formatDateUTC } from "@/utils/date";
import { debounce } from "@/utils/support";
import TransactionChange from "./transaction-change.vue";
import { ref, computed } from "vue";
import { useModalScreen } from "@/components/modal/use-modal-screen";
import { useRouter } from "vue-router";
import { useLoadingScreen } from "@/components/loading/useLoadingScreen";
import { useDialogScreen } from "@/components/dialog/use-dialog-screen";
import { useToast } from "vue-toastification";
import SummaryData from "@/components/summary-data.vue";
import { currencyBRL } from "@/components/filters/currency.filter";
import { formatCurrency } from "@/utils/numbers";

const fields = [
  { title: "Data", name: "formatted_date" },
  { title: "Descrição", name: "description" },
  { title: "Categoria", name: "category" },
  { title: "Conta", name: "account" },
  { title: { value: "Valor", clazz: "text-end" }, name: "formatted_value" },
];

const toast = useToast();
const modal = useModalScreen(TransactionChange);

const dialog = useDialogScreen(
  "Deseja realmente excluir a transação?",
  "Exclusão"
);

const router = useRouter();

let categories = [],
  accounts = [];
let transactions = ref([]);
const selectedType = ref("");
const selectedAccount = ref("");
const selectedCategory = ref(null);
const loading = useLoadingScreen();
let currentDate = new Date();

const filteredItems = computed(() =>
  transactions.value.filter((item) => {
    const typeMatch =
      !selectedType.value || item.categoryType === selectedType.value;
    const accountMatch =
      !selectedAccount.value || item.account === selectedAccount.value;
    const categoryMatch =
      !selectedCategory.value || item.categoryId === selectedCategory.value.id;

    return typeMatch && accountMatch && categoryMatch;
  })
);

// Computed property for dynamically updating the summary
const summary = computed(() =>
  filteredItems.value.reduce(
    (previous, current) => ({
      earns:
        current.categoryType === "R"
          ? previous.earns + current.value
          : previous.earns,
      expenses:
        current.categoryType === "D"
          ? previous.expenses + current.value
          : previous.expenses,
      investments:
        current.categoryType === "I"
          ? previous.investments + current.value
          : previous.investments,
    }),
    { earns: 0.0, expenses: 0.0, investments: 0.0 }
  )
);

function loadInitalData() {
  loading.show();
  const params = {
    month: currentDate.getMonth() + 1,
    year: currentDate.getFullYear(),
  };

  Promise.allSettled([
    transactionService.findAll(params),
    categoryService.findAll({ paginate: false }),
    accountService.findAll(),
  ])
    .then((results) => {
      const [respTransactions, respCategories, respAccounts] = results;

      transactions.value = mapTransactions(respTransactions.value.items);
      accounts = respAccounts.value.items;
      categories = respCategories.value.items;
    })
    .catch((err) => {
      router.push({ name: "denied" });
    })
    .finally(() => {
      loading.hide();
    });
}

function getList(month, year, filter = undefined) {
  loading.show();
  const params = {
    month: month,
    year: year,
  };
  if (filter) {
    params.search = filter;
  }

  transactionService
    .findAll(params)
    .then((resp) => {
      transactions.value = mapTransactions(resp.items);
    })
    .catch((err) => {
      console.log(err);
      router.push({ name: "denied" });
    })
    .finally(() => {
      loading.hide();
    });
}

function mapTransactions(transactionList) {
  return transactionList.map((item) => ({
    ...item,
    formatted_date: formatDateUTC(item.paymentDate, "dd/MM/yyyy"), //new Date(item.paymentDate).toLocaleDateString("pt-BR"),
    formatted_value: {
      value: currencyBRL(Math.abs(item.value)),
      clazz: item.value > 0 ? "text-success text-end" : "text-danger text-end",
    },
    value: item.value,
    category: item.category.name,
    categoryId: item.category.id,
    categoryType: item.category.type,
    account: item.account.name,
  }));
}

function exportToCSV() {
  const headers = "Data;Descrição;Valor;Categoria;Conta;Data Efetiva;Tags;Nota";

  const rows = transactions.value.map(
    (item) =>
      `${item.formatted_date};${item.description};${formatCurrency(
        "" + item.value
      )};${item.category};${item.account};${formatDateUTC(
        item.transactionDate,
        "dd/MM/yyyy"
      )};;${item.detail}`
  );

  // Combine headers and rows
  const csvContent = [
    headers, // Header row
    ...rows, // Data rows
  ].join("\n");

  // Create a Blob
  const blob = new Blob([csvContent], { type: "text/csv;charset=utf-8;" });

  // Create a download link
  const link = document.createElement("a");
  const url = URL.createObjectURL(blob);
  link.setAttribute("href", url);
  link.setAttribute("download", "table_data.csv");
  link.style.display = "none";

  // Append the link to the body and trigger download
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
}

function onTypeChange(event) {
  selectedType.value = event.target.value;
}

function onAccountChange(event) {
  selectedAccount.value = event.target.value;
}

const onChangeDebounced = debounce((newDate) => {
  currentDate = newDate;
  getList(newDate.getMonth() + 1, newDate.getFullYear());
}, 1000);

loadInitalData();

const handleEdit = async (itemClicked) => {
  const item = {
    ...itemClicked,
    value: itemClicked.value.toLocaleString("pt-BR", {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    }),
    categoryId: itemClicked.categoryId,
    transactionDate: formatDateUTC(itemClicked.transactionDate, "yyyy-MM-dd"),
    paymentDate: formatDateUTC(itemClicked.paymentDate, "yyyy-MM-dd"),
  };
  const saved = await modal.show(item);
  if (saved) {
    getList(currentDate.getMonth() + 1, currentDate.getFullYear());
  }
};

const handleDelete = async (itemClicked) => {
  const confirmed = await dialog.show();
  if (confirmed) {
    loading.show();
    transactionService
      .delete(itemClicked.id)
      .then(() => {
        toast.success("Transação Excluída com Sucesso!", {
          position: "top-center",
        });
        transactions.value = transactions.value.filter(
          (transaction) => transaction.id != itemClicked.id
        );
      })
      .finally(() => {
        loading.hide();
      });
  }
};

async function onNewClicked() {
  const saved = await modal.show();
  if (saved) {
    getList(currentDate.getMonth() + 1, currentDate.getFullYear());
  }
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
