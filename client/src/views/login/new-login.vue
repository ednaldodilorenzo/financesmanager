<template>
  <loading-screen :loading="loading" />
  <div class="authentication-card">
    <div class="card shadow rounded-0 overflow-hidden">
      <div class="row g-0">
        <div
          class="col-lg-6 bg-login d-flex align-items-center justify-content-center"
        >
          <img src="@/assets/login.png" class="img-fluid" alt="" />
        </div>
        <div class="col-lg-6">
          <div class="card-body p-4 p-sm-5">
            <h5 class="card-title">Finance Cockpit</h5>
            <p class="card-text mb-5">Você no controle das suas finanças!</p>
            <form
              :class="{
                'was-validated': v$.$dirty ? true : false,
              }"
              @submit.stop.prevent="onSubmit"
              novalidate
              class="form-body"
            >
              <div class="d-grid">
                <a class="btn btn-white radius-30" href="javascript:;"
                  ><span
                    class="d-flex justify-content-center align-items-center"
                  >
                    <img
                      class="me-2"
                      src="assets/images/icons/search.svg"
                      width="16"
                      alt=""
                    />
                    <span>Entrar com Google</span>
                  </span>
                </a>
              </div>
              <div class="login-separater text-center mb-4">
                <span>OU ENTRE COM EMAIL</span>
                <hr />
              </div>
              <div class="row g-3">
                <div class="col-12">
                  <label for="inputEmailAddress" class="form-label"
                    >Digite o Email</label
                  >
                  <div class="ms-auto position-relative">
                    <div
                      class="position-absolute top-50 translate-middle-y search-icon px-3"
                    >
                      <i class="bi bi-envelope-fill"></i>
                    </div>
                    <input
                      type="email"
                      autofocus
                      class="form-control radius-30 ps-5"
                      id="inputEmailAddress"
                      required
                      aria-required="true"
                      v-model="form.email"
                      placeholder="Email Address"
                      data-test="input-email"
                    />
                    <div class="invalid-feedback" id="live-feedback-email">
                      Formato de email inválido
                    </div>
                  </div>
                </div>
                <div class="col-12">
                  <label for="inputChoosePassword" class="form-label"
                    >Digite a Senha</label
                  >
                  <div class="ms-auto position-relative">
                    <div
                      class="position-absolute top-50 translate-middle-y search-icon px-3"
                    >
                      <i class="bi bi-lock-fill"></i>
                    </div>
                    <input
                      type="password"
                      class="form-control radius-30 ps-5"
                      id="inputChoosePassword"
                      placeholder="Digite a Senha"
                      v-model="form.password"
                      required
                    />
                    <div class="invalid-feedback" id="live-feedback-password">
                      Campo obrigatório
                    </div>
                  </div>
                </div>
                <div class="col-6">
                  <div class="form-check form-switch">
                    <input
                      class="form-check-input"
                      type="checkbox"
                      id="flexSwitchCheckChecked"
                      checked=""
                    />
                    <label class="form-check-label" for="flexSwitchCheckChecked"
                      >Lembrar</label
                    >
                  </div>
                </div>
                <div class="col-6 text-end">
                  <router-link
                    :to="{ name: ROUTE_NAMES.RECOVER }"
                    class="text-decoration-none"
                    >Esqueceu a Senha ?</router-link
                  >
                </div>
                <div class="col-12">
                  <div class="d-grid">
                    <button type="submit" class="btn btn-primary radius-30">
                      Entrar
                    </button>
                  </div>
                </div>
                <div class="col-12">
                  <p class="mb-0">
                    Ainda não possui uma conta?
                    <router-link :to="{ name: ROUTE_NAMES.SEND_MAIL }">
                      Registre-se</router-link
                    >
                  </p>
                </div>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import { useVuelidate } from "@vuelidate/core";
import { required, minLength, email } from "@vuelidate/validators";
import authService from "./auth.service";
import LoadingScreen from "@/components/loading-screen.vue";
import { ROUTE_NAMES } from "./routes.definition";
import { ROUTE_NAMES as DASHBOARD_ROUTE_NAMES } from "../dashboard/routes.definition";

export default {
  components: {
    LoadingScreen,
  },
  setup() {
    return { v$: useVuelidate() };
  },
  created() {
    this.ROUTE_NAMES = ROUTE_NAMES;
    this.DASHBOARD_ROUTE_NAMES = DASHBOARD_ROUTE_NAMES;
  },
  data() {
    return {
      form: {
        email: "",
        password: "",
        checked: [],
      },
      showInvalidLoginMessage: false,
      loading: false,
    };
  },
  validations() {
    return {
      form: {
        email: {
          required,
          email,
        },
        password: {
          required,
          minLength: minLength(3),
        },
      },
    };
  },
  methods: {
    onSubmit() {
      this.v$.$validate();

      if (this.v$.$error) {
        return;
      }

      this.loading = true;
      const { email: username, password } = this.form;
      authService
        .login(username, password)
        .then((response) => {
          if (response) {
            this.$router.push({ name: DASHBOARD_ROUTE_NAMES.INDEX });
          } else {
            this.showInvalidLoginMessage = true;
          }
        })
        .catch((e) => {
          console.log(e);

          this.showInvalidLoginMessage = true;
        })
        .finally(() => {
          this.loading = false;
        });
    },
  },
};
</script>
<style>
.authentication-card {
  display: flex;
  align-items: center;
  justify-content: center;
  margin: auto;
  max-width: 60rem;
  height: 100vh;
}

.login-separater span {
  position: relative;
  top: 26px;
  margin-top: -10px;
  background: #fff;
  padding: 5px;
  font-size: 12px;
  color: #cbcbcb;
  z-index: 1;
}
</style>
