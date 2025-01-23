const namespace = "planning";

const ROUTE_NAMES = {
  INDEX: `${namespace}-list`,
};

const ROUTES_DEFINITIONS = [
  {
    path: "/plannings",
    name: ROUTE_NAMES.INDEX,
    component: () => import("./planning-list.vue"),
    meta: {
      label: "Planejamento",
      icon: "bi bi-clipboard-check-fill",
    },
  },
];

export { ROUTE_NAMES, ROUTES_DEFINITIONS };
