import { ROUTES_DEFINITIONS as CATEGORY_ROUTES } from "./category/routes.definition";
import { ROUTES_DEFINITIONS as ACCOUNT_ROUTES } from "./account/routes.definition";
import { ROUTES_DEFINITIONS as TRANSACTION_ROUTES } from "./transaction/routes.definition";
import { ROUTES_DEFINITIONS as PLANNING_ROUTES } from "./planning/routes.definition";
import { ROUTES_DEFINITIONS as BUDGET_ROUTES } from "./budget/routes.definition";
import { ROUTES_DEFINITIONS as DASHBOARD_ROUTES } from "./dashboard/routes.definition";

const SIDEBAR_ROUTES = [
  DASHBOARD_ROUTES[0],
  CATEGORY_ROUTES[0],
  ACCOUNT_ROUTES[0],
  TRANSACTION_ROUTES[0],  
  BUDGET_ROUTES[0],
  PLANNING_ROUTES[0],
  TRANSACTION_ROUTES[1],
];
const VIEW_ROUTES = [].concat(
  DASHBOARD_ROUTES,
  CATEGORY_ROUTES,
  ACCOUNT_ROUTES,
  TRANSACTION_ROUTES,
  PLANNING_ROUTES,
  BUDGET_ROUTES
);

export { SIDEBAR_ROUTES, VIEW_ROUTES };
