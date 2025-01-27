const namespace = "account";

const ROUTE_NAMES = {
  INDEX: `${namespace}-list`,
  ADD: `${namespace}-add`,
  EDIT: `${namespace}-edit`,
};

const ROUTES_DEFINITIONS = [
  {
    path: "/accounts",
    name: ROUTE_NAMES.INDEX,
    component: () => import("./account-list.vue"),
    meta: {
      label: "Contas",
      icon: "bi bi-person-vcard",      
    },
  },
];

export { ROUTE_NAMES, ROUTES_DEFINITIONS };
