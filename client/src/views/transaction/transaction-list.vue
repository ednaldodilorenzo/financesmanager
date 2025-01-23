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
        <div class="mx-3 p-2 value-summary earn">
          Receitas {{ summary.formatted_earns }}
        </div>
        <div class="mx-3 p-2 value-summary expense">
          Despesas {{ summary.formatted_expenses }}
        </div>
        <div
          class="mx-3 p-2 value-summary"
          :class="{
            expense: summary.earns + summary.expenses < 0,
            earn: summary.earns + summary.expenses > 0,
          }"
        >
          Saldo {{ summary.formatted_balance }}
        </div>
      </div>
    </div>
  </div>
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
        :showPagination="false"
        :items="filteredItems"
        :showFilter="true"
        :actions="[
          {
            name: 'edit',
            title: 'Editar Transação',
            clazz: 'bi bi-pencil-fill',
            handler: handleEdit,
          },
          {
            name: 'delete',
            title: 'Excluir Transação',
            clazz: 'bi bi-trash-fill',
            handler: handleDelete,
          },
        ]"
        @trigger-page="getList"
        @action-clicked="onClickAction"
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
import { format } from "date-fns";
import { debounce } from "@/utils/support";
import TransactionChangeScreen from "./transaction-change-screen.vue";
import { ref, computed } from "vue";
import { useModalScreen } from "@/components/modal/use-modal-screen";
import { useRouter } from "vue-router";
import { useLoadingScreen } from "@/components/loading/useLoadingScreen";
import { useDialogScreen } from "@/components/dialog/use-dialog-screen";
import { useToast } from "vue-toastification";

const fields = [
  { title: "Data", name: "formatted_date" },
  { title: "Descrição", name: "description" },
  { title: "Categoria", name: "category" },
  { title: "Conta", name: "account" },
  { title: "Valor", name: "formatted_value" },
];

const toast = useToast();
const modal = useModalScreen(TransactionChangeScreen);

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

const filteredItems = computed(() => {
  return transactions.value.filter((item) => {
    const typeMatch =
      !selectedType.value || item.categoryType === selectedType.value;
    const accountMatch =
      !selectedAccount.value || item.account === selectedAccount.value;
    const categoryMatch =
      !selectedCategory.value || item.categoryId === selectedCategory.value.id;

    return typeMatch && accountMatch && categoryMatch;
  });
});

// Computed property for dynamically updating the summary
const summary = computed(() => {
  const summaryData = filteredItems.value.reduce(
    (previous, current) => ({
      earns:
        current.categoryType === "R"
          ? previous.earns + current.value
          : previous.earns,
      expenses:
        current.categoryType === "D"
          ? previous.expenses + current.value
          : previous.expenses,
    }),
    { earns: 0.0, expenses: 0.0 }
  );

  return {
    earns: summaryData.earns,
    expenses: summaryData.expenses,
    formatted_earns: Intl.NumberFormat("pt-br", {
      style: "currency",
      currency: "BRL",
    }).format(Math.abs(summaryData.earns)),
    formatted_expenses: Intl.NumberFormat("pt-br", {
      style: "currency",
      currency: "BRL",
    }).format(Math.abs(summaryData.expenses)),
    formatted_balance: Intl.NumberFormat("pt-br", {
      style: "currency",
      currency: "BRL",
    }).format(Math.abs(summaryData.earns + summaryData.expenses)),
  };
});

function loadInitalData() {
  const currentDate = new Date();

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
    formatted_date: format(item.paymentDate, "dd/MM/yyyy"),
    formatted_value: {
      value: Intl.NumberFormat("pt-br", {
        style: "currency",
        currency: "BRL",
      }).format(Math.abs(item.value / 100)),
      style: {
        color: item.value > 0 ? "green" : "red",
        textAlign: "right",
      },
    },
    value: item.value / 100,
    category: item.category.name,
    categoryId: item.category.id,
    categoryType: item.category.type,
    account: item.account.name,
  }));
}

function onTypeChange(event) {
  selectedType.value = event.target.value;
}

function onAccountChange(event) {
  selectedAccount.value = event.target.value;
}

const onChangeDebounced = debounce((newDate) => {
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
    transactionDate: format(itemClicked.transactionDate, "yyyy-MM-dd"),
    paymentDate: format(itemClicked.paymentDate, "yyyy-MM-dd"),
  };
  const saved = await modal.show(item);
  if (saved) {
    const currentDate = new Date();
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
    const currentDate = new Date();
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
