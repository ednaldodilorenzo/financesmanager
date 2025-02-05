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
  <summary-data
    :data="[
      {
        title: 'Total de Receitas',
        value: summary.earns,
      },
      {
        title: 'Total de Despesas',
        value: Math.abs(summary.expenses),
        percent: Math.abs(summary.expenses / summary.earns),
        percentMessage: 'do Total de Receitas',
      },
      {
        title: 'Saldo Total',
        value: summary.earns - summary.expenses,
        percent: Math.abs((summary.earns - summary.expenses) / summary.earns),
        percentMessage: 'do Total de Receitas',
      },
      {
        title: 'Investimento Total',
        value: summary.investments,
        percent: Math.abs(summary.investments / summary.earns),
        percentMessage: 'do Total de Receitas',
      },
    ]"
  />
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
            <th class="text-end">Valor</th>
            <th class="text-center">Ações</th>
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
import SummaryData from "@/components/summary-data.vue";

const loading = useLoadingScreen();
let currentDate = new Date();

const filteredItems = ref([]);
let allCategories = [];

const summary = computed(() =>
  filteredItems.value.reduce(
    (previous, current) => ({
      earns:
        current.category.type === "R"
          ? previous.earns + current.value
          : previous.earns,
      expenses:
        current.category.type === "D"
          ? previous.expenses + current.value
          : previous.expenses,
      investments:
        current.category.type === "I"
          ? previous.investments + current.value
          : previous.investments,
    }),
    { earns: 0.0, expenses: 0.0, investments: 0.0 }
  )
);

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
