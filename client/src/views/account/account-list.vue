<template>
  <div class="mt-3">
    <div class="d-flex justify-content-between">
      <h2 class="fs-4">Contas</h2>
      <nav style="--bs-breadcrumb-divider: '>'" aria-label="breadcrumb">
        <ol class="breadcrumb">
          <li class="breadcrumb-item"><a href="#">Home</a></li>
          <li class="breadcrumb-item active"><a href="#">Contas</a></li>
        </ol>
      </nav>
    </div>
  </div>
  <hr />
  <div class="d-flex flex-wrap">
    <div class="card text-center" style="width: 18rem; margin: 0.5rem">
      <div class="card-body">
        <h2>Nova</h2>
        <button
          @click="onNewClicked()"
          type="button"
          class="btn rounded-circle"
        >
          <i class="bi bi-plus-circle" style="font-size: 3rem"></i>
        </button>
      </div>
    </div>
    <account-item
      v-for="item in items"
      :key="item"
      :item="item"
      @item-edit-click="onItemEditClick"
    ></account-item>
  </div>
</template>
<script setup>
import { ref } from "vue";
import accountService from "./account.service";
import { useLoadingScreen } from "@/components/loading/useLoadingScreen";
import { useModalScreen } from "@/components/modal/use-modal-screen";
import { useRouter } from "vue-router";
import AccountChangeScreen from "./account-change-screen.vue";
import AccountItem from "./account-item.vue";

const loading = useLoadingScreen();
const items = ref([]);
const router = useRouter();
const modal = useModalScreen(AccountChangeScreen);

const getList = (filter = undefined) => {
  loading.show();
  const params = { paginate: false };
  if (filter) {
    params.search = filter;
  }

  accountService
    .findAll(params)
    .then((resp) => {
      items.value = resp.data;
    })
    .catch((err) => {
      router.push({ name: "denied" });
    })
    .finally(() => {
      loading.hide();
    });
};

getList();

const onItemEditClick = async (itemClicked) => {
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
