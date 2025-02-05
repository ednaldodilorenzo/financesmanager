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
  <summary-data
    :data="[
      {
        title: 'Planejado',
        value: summary.planned,
      },
      {
        title: 'Executado',
        value: summary.executed,
      },
      {
        title: 'Diferença',
        value: summary.planned - summary.executed,
      },
      {
        title: 'Investido',
        value: Math.abs(summary.invested),
      },
    ]"
  />
  <div class="card">
    <div class="card-body p-2">
      <bootstrap-table
        :fields="[
          { title: 'Categoria', name: 'name' },
          {
            title: 'Executado',
            name: 'formatted_value',
          },
          {
            title: { value: 'Planejado', clazz: 'text-end' },
            name: 'formatted_planned',
          },
        ]"
        :showPagination="false"
        :showNav="false"
        :items="filteredItems"
        :showFilter="false"
      >
        <template v-slot:custom-td-formatted_value="{ item, field }">
          <div class="progress" style="height: 20px">
            <div
              class="progress-bar"
              role="progressbar"
              :style="{
                width:
                  Math.abs(item.total) <= item.monthly_planned
                    ? (item.total / item.monthly_planned) * 100 + '%'
                    : '100%',
              }"
              aria-valuenow="25"
              aria-valuemin="0"
              aria-valuemax="100"
            >
              {{ currencyBRL(Math.abs(item.total)) }}
            </div>
            <div
              v-if="Math.abs(item.total) > item.monthly_planned"
              class="progress-bar bg-danger"
              role="progressbar"
              :style="{
                width:
                  (Math.abs(item.total) / item.monthly_planned - 1) * 100 + '%',
              }"
              aria-valuenow="30"
              aria-valuemin="0"
              aria-valuemax="100"
            >
              {{ currencyBRL(Math.abs(item.rest)) }}
            </div>
            <div class="d-flex justify-content-center w-75" v-else>
              {{ currencyBRL(item.rest) }}
            </div>
          </div>
        </template>
        <template v-slot:custom-td-rest="{ item, field }">
          {{ currencyBRL(item.rest) }}
        </template>
      </bootstrap-table>
    </div>
  </div>
</template>
<script setup>
import BootstrapTable from "@/components/bootstrap-table.vue";
import Calendar from "@/components/bootstrap-calendar.vue";
import planningService from "./planning.service";
import { debounce } from "@/utils/support";
import { ref, computed } from "vue";
import { useLoadingScreen } from "@/components/loading/useLoadingScreen";
import SummaryData from "@/components/summary-data.vue";
import { currencyBRL } from "@/components/filters/currency.filter";

const loading = useLoadingScreen();

const filteredItems = ref([]);

const summary = computed(() =>
  filteredItems.value.reduce(
    (previous, current) => ({
      executed: previous.executed + current.total,
      planned:
        current.type !== "D"
          ? current.type === "R"
            ? previous.planned + current.planned / 12
            : previous.planned + 0
          : previous.planned - current.planned / 12,
      invested: current.type === "I" ? previous.invested + current.total : 0,
    }),
    { executed: 0.0, planned: 0.0 }
  )
);

const getData = (month, year) => {
  loading.show();
  planningService
    .findAll({ month: month, year: year })
    .then((resp) => {
      filteredItems.value = resp.items.map((item) => ({
        ...item,
        formatted_value: {
          value: currencyBRL(Math.abs(item.total)),
          clazz:
            Math.abs(item.total) > item.planned / 12
              ? "text-danger text-end"
              : "text-success text-end",
        },
        monthly_planned: item.planned / 12,
        formatted_planned: {
          value: currencyBRL(Math.abs(item.planned / 12)),
          clazz: "text-primary text-end",
        },
        formatted_accumulated: {
          value: currencyBRL(Math.abs(item.accumulated)),
          clazz:
            Math.abs(item.accumulated) >
            (Math.abs(item.planned) / 12) * (currentDate.getMonth() + 1)
              ? "text-danger text-end"
              : "text-success text-end",
        },
        formatted_planned_accumulated: {
          value: currencyBRL(
            (item.planned / 12) * (currentDate.getMonth() + 1)
          ),
          clazz: "text-primary text-end",
        },
        formatted_total_planned: {
          value: currencyBRL(item.planned),
          clazz: "text-primary text-end",
        },
        formatted_tendency: {
          value: currencyBRL(
            Math.abs(item.accumulated) +
              ((11 - currentDate.getMonth()) * item.planned) / 12
          ),
          clazz: "text-primary text-end",
        },
        rest: item.planned / 12 - Math.abs(item.total),
      }));
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
