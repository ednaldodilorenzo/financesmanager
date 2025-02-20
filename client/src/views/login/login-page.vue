<template>
  <loading-screen :loading="loading" />
  <div class="text-center mt-2 mb-4">
    <h4 class="text-primary">Bem Vindo !</h4>
    <p class="text-muted">Entre para continuar no Finance Cockpit.</p>
  </div>
  <div
    v-if="showInvalidLoginMessage"
    class="alert alert-warning alert-dismissible fade show"
    role="alert"
    data-test="msg-invalid-login"
  >
    <strong>Login ou senha inválida!</strong>
    <button
      type="button"
      class="btn-close"
      data-bs-dismiss="alert"
      aria-label="Close"
    ></button>
  </div>
  <form
    :class="{
      'was-validated': v$.$dirty ? true : false,
    }"
    @submit.stop.prevent="onSubmit"
    novalidate
  >
    <div class="mb-3" id="groupUsername" role="group">
      <label for="inputUsername" class="form-label">Email</label>
      <input
        class="form-control"
        id="inputUsername"
        type="email"
        placeholder="Informe o email"
        required
        aria-required="true"
        v-model="form.email"
        data-test="input-email"
      />
      <div class="invalid-feedback" id="live-feedback-email">
        Formato de email inválido
      </div>
    </div>
    <div class="mb-3" id="input-group-senha" role="group">
      <div class="float-end">
        <router-link
          :to="{ name: ROUTE_NAMES.RECOVER }"
          class="text-decoration-none"
        >
          Esqueceu a Senha?</router-link
        >
      </div>
      <label for="inputPassword" class="form-label">Senha</label>
      <input
        class="form-control"
        id="inputPassword"
        type="password"
        placeholder="Informe a senha"
        required
        aria-required="true"
        v-model="form.password"
        data-test="input-password"
      />
      <div class="invalid-feedback" id="live-feedback-password">
        Campo obrigatório
      </div>
    </div>
    <div class="mb-3 form-check">
      <input type="checkbox" class="form-check-input" id="exampleCheck1" />
      <label class="form-check-label" for="exampleCheck1"
        >Lembrar meu usuário</label
      >
    </div>
    <button
      class="btn btn-primary d-block w-100 mb-3 mt-5"
      data-test="button-login"
      type="submit"
    >
      Entrar
    </button>
    <div class="mt-4 text-center">
      <p class="mb-0">
        Ainda não possui uma conta?
        <router-link
          :to="{ name: ROUTE_NAMES.SEND_MAIL }"
          class="fw-medium text-primary"
        >
          Cadastre-se agora</router-link
        >
      </p>
    </div>
  </form>
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
        food: null,
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
<style scoped>
.login-box {
  background: #fff !important;
}

input {
  background: inherit;
}
</style>
