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
  <hr />
  <div class="d-flex justify-content-center my-3">
    <Calendar @date-change="onChangeDebounced"></Calendar>
  </div>
  <!--<div class="card">
    <div class="card-body p-0">
      <div
        class="row row-cols-xxl-5 row-cols-md-3 row-cols-1 g-0 text-center align-items-center"
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
                  class="avatar-title bg-secondary-subtle fs-22 rounded-circle text-secondary"
                  ><iconify-icon
                    icon="solar:users-group-two-rounded-bold-duotone"
                  ></iconify-icon
                ></span>
              </div>
              <h3 class="mb-0 fw-bold">9,458</h3>
            </div>
            <p class="mb-0 text-muted">
              <span class="text-danger me-2"
                ><i class="ti ti-caret-down-filled"></i> 6.15%</span
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
              <h3 class="mb-0 fw-bold">R$ 10.345,00</h3>
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
              <h3 class="mb-0 fw-bold">R$ 3.000,00</h3>
            </div>
            <p class="mb-0 text-muted">
              <span class="text-success me-2"
                ><i class="ti ti-caret-up-filled"></i> 3.51%</span
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
              Pending Invoices
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
              <h3 class="mb-0 fw-bold">87.94%</h3>
            </div>
            <p class="mb-0 text-muted">
              <span class="text-danger me-2"
                ><i class="ti ti-caret-down-filled"></i> 1.05%</span
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
              Cancelled Invoices
            </h5>
            <div
              class="d-flex align-items-center justify-content-center gap-2 my-3"
            >
              <div class="avatar-sm flex-shrink-0">
                <span
                  class="avatar-title bg-danger-subtle fs-22 rounded-circle text-danger"
                  ><iconify-icon
                    icon="solar:bill-cross-bold-duotone"
                  ></iconify-icon
                ></span>
              </div>
              <h3 class="mb-0 fw-bold">7.11%</h3>
            </div>
            <p class="mb-0 text-muted">
              <span class="text-success me-2"
                ><i class="ti ti-caret-up-filled"></i> 0.05%</span
              ><span class="text-nowrap">Since last month</span>
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>-->
  <div class="card">
    <div class="card-body p-2">
      <bootstrap-table
        :fields="[
          { title: 'Categoria', name: 'name' },
          { title: 'Valor', name: 'formatted_value' },
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
import { ref } from "vue";
import { useLoadingScreen } from "@/components/loading/useLoadingScreen";

const loading = useLoadingScreen();

const filteredItems = ref([]);

const getData = (month, year) => {
  loading.show();
  planningService
    .findAll({ month: month, year: year })
    .then((resp) => {
      filteredItems.value = resp.items.map((item) => ({
        ...item,
        formatted_value: {
          value: Intl.NumberFormat("pt-br", {
            style: "currency",
            currency: "BRL",
          }).format(Math.abs(item.total / 100)),
          style: {
            color: item.total > 0 ? "green" : "red",
            textAlign: "right",
          },
        },
      }));
    })
    .finally(() => {
      loading.hide();
    });
};

const currentDate = new Date();
getData(currentDate.getMonth() + 1, currentDate.getFullYear());

const onChangeDebounced = debounce((newDate) => {
  getData(newDate.getMonth() + 1, newDate.getFullYear());
}, 1000);
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
