<template>
  <div class="text-center mt-2">
    <h4 class="text-primary">Recuperação de Senha</h4>
  </div>
  <div
    v-if="showValidationError"
    class="alert alert-warning alert-dismissible fade show"
    role="alert"
    data-test="msg-invalid-login"
  >
    <strong>{{ validationMessage }}</strong>
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
    <div class="mb-3" id="input-group-senha" role="group">
      <label for="inputPassword" class="form-label">Senha</label>
      <input
        class="form-control"
        id="inputPassword"
        type="password"
        placeholder="Enter name"
        required
        aria-required="true"
        v-model="form.password"
        data-test="input-password"
      />
      <div class="invalid-feedback" id="live-feedback-password">
        Campo obrigatório
      </div>
    </div>
    <div class="mb-3" id="input-group-confirm" role="group">
      <label for="inputConfirm" class="form-label">Confirmar senha</label>
      <input
        class="form-control"
        id="inputConfirm"
        type="password"
        placeholder="Entre a confirmação da senha"
        required
        aria-required="true"
        v-model="form.confirmPassword"
        data-test="input-confirm"
      />
      <div class="invalid-feedback">Campo obrigatório</div>
    </div>
    <button
      class="btn btn-primary d-block w-100 mb-3 mt-5"
      data-test="button-login"
      type="submit"
    >
      Submeter
    </button>
  </form>
</template>
<script>
import { useLoadingScreen } from "@/components/loading/useLoadingScreen";
import LoadingScreen from "@/components/loading-screen.vue";
import { useVuelidate } from "@vuelidate/core";
import { required, minLength, email } from "@vuelidate/validators";
import authService from "./auth.service";
import { ROUTE_NAMES } from "./routes.definition";
import { useToast } from "vue-toastification";
import { useRoute } from "vue-router";
import { HTTP_STATUS_CODE } from "@/utils/constants";

export default {
  components: {
    LoadingScreen,
  },
  setup() {
    const toast = useToast();
    const route = useRoute();

    return {
      v$: useVuelidate(),
      toast: toast,
      loading: useLoadingScreen(),
      route: route,
    };
  },
  data() {
    return {
      form: {
        password: null,
        confirmPassword: null,
      },
      showValidationError: false,
      validationMessage: "",
    };
  },
  validations() {
    return {
      form: {
        confirmPassword: {
          required,
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

      this.loading.show();
      const payload = { ...this.form };
      payload.token = this.$route.params.token;

      authService
        .redefinePassword(payload)
        .then(() => {
          this.toast.success("Senha recuperada com sucesso!", {
            position: "top-center",
          });
          this.$router.push({ name: ROUTE_NAMES.INDEX });
        })
        .catch((err) => {
          if (err.response.status === HTTP_STATUS_CODE.UNPROCESSABLE_ENTITY) {
            this.showValidationError = true;
            this.validationMessage = err.response.data.message;
          } else {
            this.toast.error("Falha no registro do usuário!");
          }
        })
        .finally(() => {
          this.loading.hide();
        });
    },
  },
};
</script>
<style scoped>
.register-box {
  background: #fff !important;
}

input {
  background: inherit;
}
</style>
