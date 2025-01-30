<template>
  <nav v-if="showNav" class="navbar bg-body-tertiary mb-3">
    <div class="d-flex w-100">
      <div>
        <a
          href="javascript:void(0)"
          v-if="showBack"
          @click="clickBack"
          class="btn btn-outline-primary me-2"
          ><i class="bx bx-arrow-back"></i
        ></a>

        <a href="#" @click="clickNew()" class="btn btn-primary">+ Novo</a>
      </div>
      <div
        v-if="showCSVButton"
        class="btn-group mx-3"
        role="group"
        aria-label="Basic outlined example"
      >
        <button type="button" @click="csvHandler" class="btn btn-success">
          <i class="bi bi-filetype-csv"></i> Exportar para CSV
        </button>
      </div>
      <div class="flex-fill">
        <input
          v-if="showFilter"
          v-autofocus
          v-model="searchQuery"
          class="form-control me-2 w-75"
          type="search"
          name="filtro"
          placeholder="Pesquisar"
          aria-label="Search"
        />
      </div>
    </div>
  </nav>
  <div v-if="filteredList.length === 0" class="data-not-found py-3">
    <h1>Nenhum dado encontrado...</h1>
  </div>
  <div v-else>
    <table
      id="tableComponent"
      class="table table-striped table-hover table-responsive"
    >
      <thead>
        <tr>
          <!-- loop through each value of the fields to get the table header -->
          <th v-for="field in fields" :class="field?.title.clazz" :key="field">
            {{
              typeof field?.title === "object"
                ? field?.title.value
                : field?.title
            }}
            <i class="bi bi-sort-alpha-down" aria-label="Sort Icon"></i>
          </th>
          <th v-if="actions.length" class="text-center">Ações</th>
        </tr>
      </thead>
      <tbody>
        <!-- Loop through the list get the each student data -->
        <tr v-for="item in filteredList" :class="item.clazz" :key="item">
          <td
            v-for="field in fields"
            :style="item[field?.name]?.style"
            :class="item[field?.name]?.clazz"
            :key="field"
          >
            {{
              typeof item[field?.name] === "object"
                ? item[field?.name].value
                : item[field?.name]
            }}
          </td>
          <td v-if="actions.length" class="text-center">
            <a
              v-for="action in actions"
              :key="action.name"
              :class="action.clazz"
              class="btn btn-sm icon-link icon-link-hover"
              data-bs-toggle="tooltip"
              data-bs-placement="top"
              :data-bs-title="action.title"
              @click.prevent="action.handler(item)"
              href="#"
            >
              <i :class="action.icon"></i>
            </a>
          </td>
        </tr>
      </tbody>
    </table>
    <nav v-if="showPagination" aria-label="Page navigation">
      <ul class="pagination justify-content-center">
        <li :class="{ disabled: pageParams?.current === 1 }" class="page-item">
          <a
            class="page-link"
            @click="setPage(1)"
            href="#"
            tabindex="-1"
            aria-disabled="true"
          >
            &lt;&lt;</a
          >
        </li>
        <li :class="{ disabled: pageParams?.current === 1 }" class="page-item">
          <a
            class="page-link"
            @click="setPage(pageParams?.prev)"
            href="#"
            tabindex="-1"
            aria-disabled="true"
          >
            &lt;</a
          >
        </li>
        <li v-for="n in pageParams?.items" :key="n" class="page-item">
          <a
            class="page-link"
            :class="{ disabled: n === pageParams?.current }"
            @click="setPage(n)"
            href="javascript:void(0)"
            >{{ n }}</a
          >
        </li>
        <li class="page-item">
          <a
            :class="{ disabled: pageParams?.current === pageParams?.max }"
            class="page-link"
            @click="setPage(pageParams.next)"
            href="#"
            >&gt;</a
          >
        </li>
        <li class="page-item">
          <a
            :class="{ disabled: pageParams?.current === pageParams?.max }"
            class="page-link"
            @click="setPage(pageParams?.max)"
            href="#"
            >&gt;&gt;</a
          >
        </li>
      </ul>
    </nav>
  </div>
</template>
<script>
import { computed, ref, watch, onMounted } from "vue";
import { debounce } from "@/utils/support";
// Importing  the lodash library
import { sortBy } from "lodash";

export default {
  name: "TableComponent",
  emits: [
    "trigger-page",
    "new-clicked",
    "action-clicked",
    "search-input",
    "back-clicked",
  ],
  props: {
    items: {
      type: [Array, Function],
      default: () => [],
    },
    csvHandler: {
      type: Function,
      default: () => {},
    },
    showCSVButton: {
      type: Boolean,
      default: false,
    },
    pageSize: {
      type: Number,
      required: true,
      default: 10,
    },
    totalPages: {
      type: Number,
      required: true,
      default: 0,
    },
    currentPage: {
      type: Number,
      required: true,
      default: 1,
    },
    fields: {
      type: Array,
    },
    actions: {
      type: Array,
      default: () => [],
    },
    showFilter: {
      type: Boolean,
      default: () => true,
    },
    showBack: {
      type: Boolean,
      default: () => false,
    },
    showPagination: {
      type: Boolean,
      default: () => true,
    },
    showNav: {
      type: Boolean,
      default: () => true,
    },
  },
  methods: {
    clickNew() {
      this.$emit("new-clicked");
    },
    clickBack() {
      this.$emit("back-clicked");
    },
  },
  setup(props) {
    let searchQuery = ref("");
    let totalRecords = ref(0);
    let currentPage = ref(1);
    const filteredItems = ref([]);

    const fetchItems = (page, limit, query) => {
      if (typeof props.items === "function") {
        return props.items(page, limit, query).then((result) => result);
      } else {
        return new Promise((resolve, reject) => {
          let result = props.items;
          let filteredItems = [];

          if (query) {
            const filter = query.toLowerCase();
            filteredItems = result.filter((item) =>
              Object.values(item).some((value) =>
                String(value).toLowerCase().includes(filter)
              )
            );
          } else {
            filteredItems = props.items;
          }

          const firstIndex = page * limit - limit;
          const lastIndex = firstIndex + (limit - 1);

          return resolve({
            items: props.showPagination
              ? filteredItems.slice(firstIndex, lastIndex)
              : filteredItems,
            total: filteredItems.length,
            page: page,
          });
        });
      }
    };

    const loadData = (page, limit, query) => {
      fetchItems(page, limit, query).then((result) => {
        filteredItems.value = result.items;
        totalRecords.value = result.total;
        currentPage.value = result.page;
      });
    };

    // Debounced function for calling `props.items`
    const fetchItemsDebounced = debounce(loadData, 1000); // Debounce for 1s

    onMounted(() => {
      loadData(currentPage.value, props.pageSize, "");
    });

    const setPage = async (pageNumber) => {
      loadData(pageNumber, props.pageSize, searchQuery.value);
    };

    watch(searchQuery, () => {
      fetchItemsDebounced(currentPage.value, props.pageSize, searchQuery.value);
    });

    watch(
      () => props.items,
      () => {
        loadData(currentPage.value, props.pageSize, "");
      }
    );

    // Computed property for accessing filtered items
    const filteredList = computed(() => {
      return filteredItems.value;
    });

    const pageParams = computed(() => {
      const current = currentPage.value;
      const max = Math.ceil(totalRecords.value / 10);

      if (!current || !max) {
        return null;
      }

      let prev = current === 1 ? null : current - 1,
        next = current === max ? null : current + 1,
        items = [1];

      if (current === 1 && max === 1) return { current, prev, next, items };
      if (current > 4) items.push("…");

      let r = 2,
        r1 = current - r,
        r2 = current + r;

      for (let i = r1 > 2 ? r1 : 2; i <= Math.min(max, r2); i++) items.push(i);

      if (r2 + 1 < max) items.push("…");
      if (r2 < max) items.push(max);

      return { current, prev, next, items, max };
    });

    return { searchQuery, filteredList, pageParams, setPage };
  },
};
</script>
<style scoped>
table th:hover {
  background: #f2f2f2;
}
</style>
