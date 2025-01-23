<template>
  <div class="mt-3">
    <div class="d-flex justify-content-between">
      <h2 class="fs-4">Orçamento</h2>
      <nav style="--bs-breadcrumb-divider: '>'" aria-label="breadcrumb">
        <ol class="breadcrumb">
          <li class="breadcrumb-item"><a href="#">Home</a></li>
          <li class="breadcrumb-item active"><a href="#">Orçamento</a></li>
        </ol>
      </nav>
    </div>
  </div>
  <div class="card mb-3">
    <div class="card-body p-2">
      <div class="d-flex justify-content-center my-3">
        <Calendar
          @date-change="onChangeDebounced"
          :only-years="true"
        ></Calendar>
      </div>
    </div>
  </div>
  <div class="card mb-3">
    <div class="card-body p-0">
      <div
        class="row row-cols-xxl-5 row-cols-md-3 row-cols-1 g-0 text-center align-items-center justify-content-center"
      >
        <div class="col border-end border-light border-dashed">
          <div class="mt-3 mt-md-0 p-3">
            <h5
              title="Number of Orders"
              class="text-muted fs-13 text-uppercase"
            >
              Total de Receitas
            </h5>
            <div
              class="d-flex align-items-center justify-content-center gap-2 my-3"
            >
              <div class="avatar-sm flex-shrink-0">
                <span
                  class="avatar-title bg-primary-subtle fs-22 rounded-circle text-primary"
                  ><iconify-icon
                    icon="solar:bill-list-bold-duotone"
                  ></iconify-icon
                ></span>
              </div>
              <h3 class="mb-0 fw-bold">{{ summary.formatted_earns }}</h3>
            </div>
            <p class="mb-0 text-muted">
              <span class="text-success me-2"
                ><i class="ti ti-caret-up-filled"></i> 26.87%</span
              ><span class="text-nowrap">Since last month</span>
            </p>
          </div>
        </div>
        <div class="col border-end border-light border-dashed">
          <div class="mt-3 mt-md-0 p-3">
            <h5
              title="Number of Orders"
              class="text-muted fs-13 text-uppercase"
            >
              Total de Despesas
            </h5>
            <div
              class="d-flex align-items-center justify-content-center gap-2 my-3"
            >
              <div class="avatar-sm flex-shrink-0">
                <span
                  class="avatar-title bg-warning-subtle fs-22 rounded-circle text-warning"
                  ><iconify-icon
                    icon="solar:wallet-money-bold-duotone"
                  ></iconify-icon
                ></span>
              </div>
              <h3 class="mb-0 fw-bold">{{ summary.formatted_expenses }}</h3>
            </div>
            <p class="mb-0 text-muted">
              <span class="text-success me-2"
                ><i class="ti ti-caret-up-filled"></i>
                {{ summary.formatted_expense_percentage }}</span
              ><span class="text-nowrap">do Total de Receitas</span>
            </p>
          </div>
        </div>
        <div class="col border-end border-light border-dashed">
          <div class="mt-3 mt-md-0 p-3">
            <h5
              title="Number of Orders"
              class="text-muted fs-13 text-uppercase"
            >
              Saldo Total
            </h5>
            <div
              class="d-flex align-items-center justify-content-center gap-2 my-3"
            >
              <div class="avatar-sm flex-shrink-0">
                <span
                  class="avatar-title bg-success-subtle fs-22 rounded-circle text-success"
                  ><iconify-icon
                    icon="solar:banknote-2-bold-duotone"
                  ></iconify-icon
                ></span>
              </div>
              <h3 class="mb-0 fw-bold">{{ summary.formatted_balance }}</h3>
            </div>
            <p class="mb-0 text-muted">
              <span class="text-danger me-2"
                ><i class="ti ti-caret-down-filled"></i> 1.05%</span
              ><span class="text-nowrap">Since last month</span>
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
  <div class="card">
    <div class="card-body p-2">
      <nav class="navbar bg-body-tertiary mb-3">
        <div class="d-flex w-100">
          <div style="width: 15%">
            <a href="#" @click="clickNew()" class="btn btn-primary">+ Novo</a>
          </div>
        </div>
      </nav>
      <table class="table table-striped table-hover table-responsive">
        <thead>
          <tr>
            <!-- loop through each value of the fields to get the table header -->
            <th>Categoria</th>
            <th>Valor</th>
            <th>Ações</th>
          </tr>
        </thead>
        <tbody>
          <budget-record
            v-for="item of filteredItems"
            :budget-item="item"
            :options="allCategories"
            @cancel-item="cancelRecord"
            @delete-item="handleDelete"
            @change-item="handleChange"
          ></budget-record>
        </tbody>
      </table>
    </div>
  </div>
</template>
<script setup>
import Calendar from "@/components/bootstrap-calendar.vue";
import budgetService from "./budget.service";
import categoryService from "@/views/category/category.service";
import { debounce } from "@/utils/support";
import { ref, computed } from "vue";
import { useLoadingScreen } from "@/components/loading/useLoadingScreen";
import BudgetRecord from "./budget-record.vue";
import BudgetItem from "./budget-item";

const loading = useLoadingScreen();
let currentDate = new Date();

const filteredItems = ref([]);
let allCategories = [];

const summary = computed(() => {
  const summaryData = filteredItems.value.reduce(
    (previous, current) => ({
      earns:
        current.category.type === "R"
          ? previous.earns + current.value
          : previous.earns,
      expenses:
        current.category.type === "D"
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
    }).format(Math.abs(summaryData.earns - summaryData.expenses)),
    formatted_expense_percentage: Intl.NumberFormat("pt-br", {
      style: "percent",
      currency: "BRL",
    }).format(summaryData.expenses / summaryData.earns),
  };
});

const handleDelete = (item) => {
  filteredItems.value = filteredItems.value.filter(
    (value) => value.id !== item.id
  );
};

const cancelRecord = (item) => {
  filteredItems.value = filteredItems.value.filter((value) => value !== item);
};

const handleChange = (updatedItem) => {
  const index = filteredItems.value.findIndex(
    (item) => item.id === updatedItem.id
  );

  if (index !== -1) filteredItems.value[index] = updatedItem;
};

const getData = () => {
  loading.show();
  Promise.allSettled([
    categoryService.findAll({ paginate: false }),
    budgetService.findAll({ year: currentDate.getFullYear() }),
  ])
    .then((results) => {
      const [categoryResults, budgetResults] = results;
      allCategories = categoryResults.value.items;
      parseBudgetResponse(budgetResults.value);
    })
    .finally(() => {
      loading.hide();
    });
};

const clickNew = () => {
  filteredItems.value.push(
    new BudgetItem(
      {
        year: currentDate.getFullYear(),
        value: 0,
        categoryId: 0,
      },
      allCategories
    )
  );
};

const parseBudgetResponse = (resp) => {
  filteredItems.value = resp.items.map(
    (item) => new BudgetItem(item, allCategories)
  );
};

const getBudget = (year) => {
  loading.show();
  budgetService
    .findAll({ year: year })
    .then(parseBudgetResponse)
    .finally(() => {
      loading.hide();
    });
};

getData();

const onChangeDebounced = debounce((newDate) => {
  currentDate = newDate;
  getBudget(currentDate.getFullYear());
}, 1000);
</script>
