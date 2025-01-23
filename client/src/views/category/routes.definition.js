const namespace = "category";

const ROUTE_NAMES = {
  INDEX: `${namespace}-list`,
  ADD: `${namespace}-add`,
  EDIT: `${namespace}-edit`,
};

const ROUTES_DEFINITIONS = [
  {
    path: "/categories",
    name: ROUTE_NAMES.INDEX,
    component: () => import("./category-list.vue"),
    meta: {
      label: "Categorias",
      icon: "bi bi-bookmark-plus-fill",
    },
  },
];

export { ROUTE_NAMES, ROUTES_DEFINITIONS };
