const namespace = "auth";

const ROUTE_NAMES = {
  INDEX: `${namespace}`,
  LOGIN: `${namespace}-login`,
  REGISTER: `${namespace}-register`,
};

const ROUTES_DEFINITIONS = [
  {
    path: "/auth",
    name: ROUTE_NAMES.INDEX,
    component: () => import("./auth-base.vue"),
    children: [
      {
        path: "/login",
        name: ROUTE_NAMES.LOGIN,
        component: () => import("./login-page.vue"),
      },
      {
        path: "/register",
        name: ROUTE_NAMES.REGISTER,
        component: () => import("./registration-page.vue"),
      },
    ],
  },
];

export { ROUTE_NAMES, ROUTES_DEFINITIONS };
