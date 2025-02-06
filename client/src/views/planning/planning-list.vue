<template>
  <div class="mt-3">
    <div class="d-flex justify-content-between">
      <h2 class="fs-4">Planejamento</h2>
      <nav style="--bs-breadcrumb-divider: '>'" aria-label="breadcrumb">
        <ol class="breadcrumb">
          <li class="breadcrumb-item"><a href="#">Home</a></li>
          <li class="breadcrumb-item active"><a href="#">Planejamento</a></li>
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
  <div class="card">
    <div class="card-body p-2">
      <bootstrap-table
        :fields="[
          { title: 'Receitas', name: 'name' },
          {
            title: '',
            name: 'formatted_value',
          },
          {
            title: { value: 'Planejado', clazz: 'text-end' },
            name: 'formatted_planned',
          },
        ]"
        :showPagination="false"
        :showNav="false"
        :items="earnsList"
        :showFilter="false"
      >
        <template #first-row>
          <tr>
            <td>Total</td>
            <td>
              <bootstrap-plan-exec-bar
                :planned="earnsSummary.planned"
                :executed="earnsSummary.executed"
              />
            </td>
            <td class="text-end text-primary">
              {{ currencyBRL(earnsSummary.planned) }}
            </td>
          </tr>
        </template>
        <template #custom-td-formatted_value="{ item, field }">
          <bootstrap-plan-exec-bar
            :planned="item.monthly_planned"
            :executed="Math.abs(item.total)"
          />
        </template>
      </bootstrap-table>
      <bootstrap-table
        :fields="[
          { title: 'Despesas', name: 'name' },
          {
            title: '',
            name: 'chartValues',
          },
          {
            title: { value: 'Planejado', clazz: 'text-end' },
            name: 'formatted_planned',
          },
        ]"
        :showPagination="false"
        :showNav="false"
        :items="expensesList"
        :showFilter="false"
      >
        <template #first-row>
          <tr>
            <td>Total</td>
            <td>
              <bootstrap-plan-exec-bar
                :planned="expensesSummary.planned"
                :executed="expensesSummary.executed"
              />
            </td>
            <td class="text-end text-primary">
              {{ currencyBRL(expensesSummary.planned) }}
            </td>
          </tr>
        </template>
        <template #custom-td-chartValues="{ item, field }">
          <bootstrap-plan-exec-bar
            :planned="item.monthly_planned"
            :executed="Math.abs(item.total)"
          />
        </template>
      </bootstrap-table>
    </div>
  </div>
</template>
<script setup>
import BootstrapTable from "@/components/bootstrap-table.vue";
import BootstrapPlanExecBar from "@/components/boostra-planexec-bar.vue";
import Calendar from "@/components/bootstrap-calendar.vue";
import planningService from "./planning.service";
import { debounce } from "@/utils/support";
import { ref, computed } from "vue";
import { useLoadingScreen } from "@/components/loading/useLoadingScreen";
import { currencyBRL } from "@/components/filters/currency.filter";

const loading = useLoadingScreen();

const expensesList = ref([]);
const earnsList = ref([]);

const earnsSummary = computed(() =>
  earnsList.value.reduce(
    (previous, current) => ({
      planned: previous.planned + current.monthly_planned,
      executed: previous.executed + current.total,
    }),
    { planned: 0.0, executed: 0.0 }
  )
);

const expensesSummary = computed(() =>
  expensesList.value.reduce(
    (previous, current) => ({
      planned: previous.planned + current.monthly_planned,
      executed: previous.executed + Math.abs(current.total),
    }),
    { planned: 0.0, executed: 0.0 }
  )
);

const getData = (month, year) => {
  loading.show();
  planningService
    .findAll({ month: month, year: year })
    .then((resp) => {
      const respList = resp.items.map((item) => ({
        ...item,
        monthly_planned: item.planned / 12,
        formatted_planned: {
          value: currencyBRL(Math.abs(item.planned / 12)),
          clazz: "text-primary text-end",
        },
        rest: item.planned / 12 - Math.abs(item.total),
      }));

      expensesList.value = respList.filter((item) => item.type === "D");
      earnsList.value = respList.filter((item) => item.type === "R");
    })
    .finally(() => {
      loading.hide();
    });
};

let currentDate = new Date();
getData(currentDate.getMonth() + 1, currentDate.getFullYear());

const onChangeDebounced = debounce((newDate) => {
  currentDate = newDate;
  getData(newDate.getMonth() + 1, newDate.getFullYear());
}, 1000);
</script>
