import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import "bootstrap/dist/css/bootstrap.min.css";
import "bootstrap";
import "bootstrap-icons/font/bootstrap-icons.css";
import "./assets/styles.css";
import Toast from "vue-toastification";
import "vue-toastification/dist/index.css";
import currencyDirective from "./components/currency.directive";

const app = createApp(App);

app.use(Toast, {
  transition: "Vue-Toastification__bounce",
  maxToasts: 20,
  newestOnTop: true,
});

app.use(store);
app.use(router);

app.directive("currency", currencyDirective);

app.mount("#app");
