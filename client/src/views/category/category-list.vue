<template>
  <div class="mt-3">
    <div class="d-flex justify-content-between">
      <h2 class="fs-4">Categorias</h2>
      <nav style="--bs-breadcrumb-divider: '>'" aria-label="breadcrumb">
        <ol class="breadcrumb">
          <li class="breadcrumb-item"><a href="#">Home</a></li>
          <li class="breadcrumb-item active"><a href="#">Categorias</a></li>
        </ol>
      </nav>
    </div>
  </div>
  <hr />
  <Table
    :fields="[
      { title: 'Nome', name: 'name' },
      { title: 'Tipo', name: 'type_desc' },
    ]"
    :showPagination="true"
    :items="items"
    :actions="[
      {
        name: 'edit',
        title: 'Editar Categoria',
        clazz: 'bi bi-pencil-fill',
        handler: handleEdit,
      },
    ]"
    :showFilter="true"
    @search-input="onChangeDebounced"
    @action-clicked="onClickAction"
    @new-clicked="onNewClicked()"
  ></Table>
</template>
<script setup>
// Importing the table component
import { ref } from "vue";
import Table from "@/components/bootstrap-table.vue";
import categoryService from "./category.service";
import { debounce } from "@/utils/support";
import { useLoadingScreen } from "@/components/loading/useLoadingScreen";
import CategoryChangeScreen from "./category-change-screen.vue";
import { useRouter } from "vue-router";
import { useModalScreen } from "@/components/modal/use-modal-screen";

const loading = useLoadingScreen();
const router = useRouter();
const items = ref([]);

const modal = useModalScreen(CategoryChangeScreen);

const getList = async () => {
  loading.show();
  const params = { paginate: false };

  try {
    const resp = await categoryService.findAll(params);
    items.value = resp.items.map((item) => ({
      ...item,
      type_desc: item.type === "D" ? "Despesa" : "Receita",
    }));
  } catch (err) {
    console.log(err);
    router.push({ name: "denied" });
  } finally {
    loading.hide();
  }
};

getList();

// const getList = async (page = 1, pageSize = 10, filter = undefined) => {
//   loading.show();
//   const params = { paginate: true, page: page, pageSize: pageSize };
//   if (filter) {
//     params.filter = filter;
//   }

//   try {
//     const resp = await categoryService.findAll(params);

//     return {
//       ...resp,
//       items: resp.items.map((item) => ({
//         ...item,
//         type_desc: item.type === "D" ? "Despesa" : "Receita",
//       })),
//     };
//   } catch (err) {
//     console.log(err);
//     router.push({ name: "denied" });
//   } finally {
//     loading.hide();
//   }
// };

const onChangeDebounced = debounce((event) => {
  this.getList(event.target.value);
}, 1000);

const handleEdit = async (itemClicked) => {
  console.log(JSON.stringify(itemClicked));
  const saved = await modal.show(itemClicked);
  if (saved) {
    getList();
  }
};

const onNewClicked = async () => {
  const saved = await modal.show();
  if (saved) {
    getList();
  }
};
</script>
