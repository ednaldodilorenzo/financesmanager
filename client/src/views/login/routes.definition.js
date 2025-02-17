const namespace = "auth";

const ROUTE_NAMES = {
  INDEX: `${namespace}`,
  LOGIN: `${namespace}-login`,
  REGISTER: `${namespace}-register`,
  VERIFY: `${namespace}-verify`,
  SEND_MAIL: `${namespace}-send-mail`,
  REDEFINE: `${namespace}-redefine`,
  RECOVER: `${namespace}-recover`,
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
        path: "/register/:token",
        name: ROUTE_NAMES.REGISTER,
        component: () => import("./registration-page.vue"),
      },
      {
        path: "/verify/:token",
        name: ROUTE_NAMES.VERIFY,
        component: () => import("./confirmation-page.vue"),
      },
      {
        path: "/registration",
        name: ROUTE_NAMES.SEND_MAIL,
        component: () => import("./send-mail.vue"),
      },
      {
        path: "/redefine/:token",
        name: ROUTE_NAMES.REDEFINE,
        component: () => import("./redefine-password.vue"),
      },
      {
        path: "/recover",
        name: ROUTE_NAMES.RECOVER,
        component: () => import("./send-mail.vue"),
      },
    ],
  },
];

export { ROUTE_NAMES, ROUTES_DEFINITIONS };
