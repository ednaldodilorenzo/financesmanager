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
        title: 'Executado',
        value: summary.executed,
      },
      {
        title: 'Planejado',
        value: summary.planned,
      },
      {
        title: 'Diferença',
        value: summary.planned - summary.executed,
      },
    ]"
  />
  <div class="card">
    <div class="card-body p-2">
      <bootstrap-table
        :fields="[
          { title: 'Categoria', name: 'name' },
          { title: 'Executado', name: 'formatted_value' },
          { title: 'Planejado', name: 'formatted_planned' },
          { title: 'Executado Acumulado', name: 'formatted_accumulated' },
          {
            title: 'Planejado Acumulado',
            name: 'formatted_planned_accumulated',
          },
        ]"
        :showPagination="false"
        :showNav="false"
        :items="filteredItems"
        :showFilter="false"
        :actions="[]"
      ></bootstrap-table>
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

const summary = computed(() => {
  return filteredItems.value.reduce(
    (previous, current) => ({
      executed: Math.abs(previous.executed + current.total),
      planned: previous.planned + current.planned / 12,
    }),
    { executed: 0.0, planned: 0.0 }
  );
});

const getData = (month, year) => {
  loading.show();
  planningService
    .findAll({ month: month, year: year })
    .then((resp) => {
      filteredItems.value = resp.items.map((item) => ({
        ...item,
        formatted_value: {
          value: currencyBRL(Math.abs(item.total)),
          style: {
            textAlign: "center",
          },
          clazz:
            Math.abs(item.total) > item.planned / 12
              ? "text-danger"
              : "text-success",
        },
        formatted_planned: {
          value: currencyBRL(Math.abs(item.planned / 12)),
          style: {
            textAlign: "center",
          },
          clazz: "text-primary",
        },
        formatted_accumulated: {
          value: currencyBRL(Math.abs(item.accumulated)),
          style: {
            textAlign: "center",
          },
          clazz:
            Math.abs(item.accumulated) >
            (Math.abs(item.planned) / 12) * (currentDate.getMonth() + 1)
              ? "text-danger"
              : "text-success",
        },
        formatted_planned_accumulated: {
          value: currencyBRL(
            (Math.abs(item.planned) / 12) * (currentDate.getMonth() + 1)
          ),
          style: {
            textAlign: "center",
          },
          clazz: "text-primary",
        },
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
