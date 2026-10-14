const namespace = "transaction";

const ROUTE_NAMES = {
  INDEX: `${namespace}-list`,
  ADD: `${namespace}-add`,
  EDIT: `${namespace}-edit`,
  UPLOAD: `${namespace}-upload`,
  REPORT: `${namespace}-report`
};

const ROUTES_DEFINITIONS = [
  {
    path: "/transactions",
    name: ROUTE_NAMES.INDEX,
    component: () => import("./transaction-list.vue"),
    meta: {
      label: "Transações",
      icon: "bi bi-cash-coin",
    },
  },
  {
    path: "/import-transactions",
    name: ROUTE_NAMES.UPLOAD,
    component: () => import("./transaction-import.vue"),
    meta: {
      label: "Importação",
      icon: "bi bi-file-arrow-up-fill",
    },
  },
  {
    path: "/transaction-reports",
    name: ROUTE_NAMES.REPORT,
    component: () => import("./transaction-report.vue"),
    meta: {
      label: "Relatórios",
      icon: "bi bi-file-earmark-text-fill",
    },
  }
];

export { ROUTE_NAMES, ROUTES_DEFINITIONS };
