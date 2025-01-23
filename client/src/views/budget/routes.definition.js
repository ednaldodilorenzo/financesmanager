const namespace = "budget";

const ROUTE_NAMES = {
  INDEX: `${namespace}-list`,
};

const ROUTES_DEFINITIONS = [
  {
    path: "/budgets",
    name: ROUTE_NAMES.INDEX,
    component: () => import("./budget-list.vue"),
    meta: {
      label: "Orçamento",
      icon: "bi bi-clipboard-check-fill",
    },
  },
];

export { ROUTE_NAMES, ROUTES_DEFINITIONS };
