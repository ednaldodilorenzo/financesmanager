import { createRouter, createWebHistory } from "vue-router";
import {
  ROUTES_DEFINITIONS as LOGIN_ROUTES,
  ROUTE_NAMES as LOGIN_ROUTE_NAMES,
} from "@/views/login/routes.definition";
import { VIEW_ROUTES } from "@/views/routes.definition";
import Cookies from "js-cookie";

const routes = [
  {
    path: "/",
    name: "admin",
    component: () => import("@/views/admin/admin-layout.vue"),
    children: [
      {
        path: "/denied",
        name: "denied",
        component: () => import("@/views/error/access-denied.vue"),
      },
      ...VIEW_ROUTES,
    ],
  },
  ...LOGIN_ROUTES,
];

const PUBLIC_ROUTES_NAMES = [
  LOGIN_ROUTE_NAMES.LOGIN,
  LOGIN_ROUTE_NAMES.REGISTER,
  LOGIN_ROUTE_NAMES.SEND_MAIL,
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

router.beforeEach((to, from, next) => {
  if (!PUBLIC_ROUTES_NAMES.includes(to.name) && !Cookies.get("jwtToken")) {
    next({ name: LOGIN_ROUTE_NAMES.LOGIN });
  } else {
    next();
  }
});

export default router;
