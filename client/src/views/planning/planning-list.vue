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
      <div class="d-flex flex-column align-items-center my-3">
        <Calendar class="mb-3" @date-change="onChangeDebounced"></Calendar>
        <bootstrap-select
          style="text-align: -webkit-center;"
          @change="updateState"
          v-model="type"
          class="w-50"
          :options="[
            { id: 'M', value: 'Mensal' },
            { id: 'Y', value: 'Anual' },
          ]"
          :key-field="'id'"
          :value-field="'value'"
        >
        </bootstrap-select>
      </div>
    </div>
  </div>
  <div class="card">
    <div class="card-body p-2">
      <bootstrap-table
        :fields="[
          { title: 'Investimentos', name: 'name' },
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
        :items="investmentsList"
        :showFilter="false"
      >
        <template #first-row>
          <tr>
            <td class="w-25">Total</td>
            <td class="w-50">
              <bootstrap-plan-exec-bar
                :planned="investmentsSumary.planned"
                :executed="investmentsSumary.executed"
                :percent-divider="percentDivider"
              />
            </td>
            <td class="text-end text-primary w-25">
              {{ currencyBRL(investmentsSumary.planned) }}
            </td>
          </tr>
        </template>
        <template #custom-td-formatted_value="{ item, field }">
          <td>
            <bootstrap-plan-exec-bar
              :planned="item.planned"
              :executed="item.executed"
              :percent-divider="percentDivider"
            />
          </td>
        </template>
      </bootstrap-table>
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
            <td class="w-25">Total</td>
            <td class="w-50">
              <bootstrap-plan-exec-bar
                :planned="earnsSummary.planned"
                :executed="earnsSummary.executed"
                :percent-divider="percentDivider"
              />
            </td>
            <td class="text-end text-primary w-25">
              {{ currencyBRL(earnsSummary.planned) }}
            </td>
          </tr>
        </template>
        <template #custom-td-formatted_value="{ item, field }">
          <td>
            <bootstrap-plan-exec-bar
              :planned="item.planned"
              :executed="item.executed"
              :percent-divider="percentDivider"
            />
          </td>
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
            <td class="w-25">Total</td>
            <td class="w-50">
              <bootstrap-plan-exec-bar
                :planned="expensesSummary.planned"
                :executed="expensesSummary.executed"
                :percent-divider="percentDivider"
              />
            </td>
            <td class="text-end text-primary w-25">
              {{ currencyBRL(expensesSummary.planned) }}
            </td>
          </tr>
        </template>
        <template #custom-td-chartValues="{ item, field }">
          <td>
            <bootstrap-plan-exec-bar
              :planned="item.planned"
              :executed="item.executed"
              :percent-divider="percentDivider"
            />
          </td>
        </template>
      </bootstrap-table>
    </div>
  </div>
</template>
<script setup>
import BootstrapTable from "@/components/bootstrap-table.vue";
import BootstrapPlanExecBar from "@/components/bootstrap-planexec-bar.vue";
import BootstrapSelect from "@/components/bootstrap-select.vue";
import Calendar from "@/components/bootstrap-calendar.vue";
import planningService from "./planning.service";
import { debounce } from "@/utils/support";
import { ref, computed } from "vue";
import { useLoadingScreen } from "@/components/loading/useLoadingScreen";
import { currencyBRL } from "@/components/filters/currency.filter";

const loading = useLoadingScreen();

const expensesList = ref([]);
const earnsList = ref([]);
const investmentsList = ref([]);
let fullList = [];
const type = ref("M");
const percentDivider = ref(0);

const investmentsSumary = computed(() => 
  investmentsList.value.reduce(
    (previous, current) => ({
      planned: previous.planned + current.planned,
      executed: previous.executed + current.executed,
    }),
    { planned: 0.0, executed: 0.0 }
  )
);

const earnsSummary = computed(() =>
  earnsList.value.reduce(
    (previous, current) => ({
      planned: previous.planned + current.planned,
      executed: previous.executed + current.executed,
    }),
    { planned: 0.0, executed: 0.0 }
  )
);

const expensesSummary = computed(() =>
  expensesList.value.reduce(
    (previous, current) => ({
      planned: previous.planned + current.planned,
      executed: previous.executed + current.executed,
    }),
    { planned: 0.0, executed: 0.0 }
  )
);

const updateState = () => {
  const respList =
    type.value === "M"
      ? fullList.map((item) => ({
          ...item,
          planned: item.planned / 12,
          executed: Math.abs(item.total),
          formatted_planned: {
            value: currencyBRL(Math.abs(item.planned / 12)),
            clazz: "text-primary text-end",
          },
        }))
      : fullList.map((item) => ({
          ...item,
          executed: Math.abs(item.accumulated),
          formatted_planned: {
            value: currencyBRL(Math.abs(item.planned)),
            clazz: "text-primary text-end",
          },
        }));

  percentDivider.value =
    type.value === "M" ? 0 : (currentDate.getMonth() + 1) / 12;

  expensesList.value = respList.filter((item) => item.type === "D");
  earnsList.value = respList.filter((item) => item.type === "R");
  investmentsList.value = respList.filter((item) => item.type === "I");
};

const getData = (month, year) => {
  loading.show();
  planningService
    .findAll({ month: month, year: year })
    .then((resp) => {
      fullList = resp.data;
      updateState();
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
