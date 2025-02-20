const namespace = "dashboard";

const ROUTE_NAMES = {
  INDEX: `${namespace}`,
};

const ROUTES_DEFINITIONS = [
  {
    path: "/dashboard",
    name: ROUTE_NAMES.INDEX,
    component: () => import("./dashboard-main.vue"),
    meta: {
      label: "Dashboard",
      icon: "bi bi-house-door-fill me-2",
    },
  },
];

export { ROUTE_NAMES, ROUTES_DEFINITIONS };
