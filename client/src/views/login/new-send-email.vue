<template>
  <div class="container-fluid">
    <div class="authentication-card">
      <div class="card shadow rounded-0 overflow-hidden">
        <div class="row g-0">
          <div
            class="col-lg-6 d-flex align-items-center justify-content-center border-end"
          >
            <img
              src="@/assets/forgot-password-frent-img.jpg"
              class="img-fluid"
              alt=""
            />
          </div>
          <div class="col-lg-6">
            <div class="card-body p-4 p-sm-5">
              <h5 class="card-title">Esqueceu a Senha?</h5>
              <p class="card-text mb-5">
                Informe o email registrado para reiniciar a senha
              </p>
              <form
                @submit.prevent="sendEmail"
                class="form-body"
                autocomplete="off"
              >
                <div class="row g-3">
                  <div class="col-12">
                    <label for="inputEmailid" class="form-label">Email</label>
                    {{ invalidEmailMessage }}
                    <input
                      type="email"
                      v-model="email"
                      autofocus
                      class="form-control form-control-lg radius-30"
                      :class="{
                        'is-invalid': invalidEmailMessage.trim() !== '',
                      }"
                      id="inputEmailid"
                      placeholder="Email"
                    />
                    <div class="invalid-feedback" id="live-feedback-password">
                      {{ invalidEmailMessage }}
                    </div>
                  </div>
                  <div class="col-12">
                    <div class="d-grid gap-3">
                      <button
                        type="submit"
                        class="btn btn-lg btn-primary radius-30"
                      >
                        Enviar
                      </button>
                      <router-link
                        class="btn btn-lg btn-light radius-30"
                        :to="{ name: ROUTE_NAMES.LOGIN }"
                        >Voltar para o Login</router-link
                      >
                    </div>
                  </div>
                </div>
              </form>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup>
import authService from "./auth.service";
import { ref } from "vue";
import { useLoadingScreen } from "@/components/loading/useLoadingScreen";
import { HTTP_STATUS_CODE } from "@/utils/constants";
import { useToast } from "vue-toastification";
import { useRoute } from "vue-router";
import { ROUTE_NAMES } from "./routes.definition";

const email = ref("");
const invalidEmailMessage = ref("");
const loading = useLoadingScreen();
const toast = useToast();
const route = useRoute();

const showValidationError = ref(false);

const sendEmail = () => {
  if (!email.value.trim()) {
    invalidEmailMessage.value = "Campo Obrigatório";
    return;
  }
  loading.show();
  const method =
    route.name === ROUTE_NAMES.RECOVER
      ? authService.startRecoverProcess(email.value)
      : authService.startRegistration(email.value);
  method
    .then(() => {
      toast.success("Email enviado com sucesso!");
      invalidEmailMessage.value = "";
      email.value = "";
    })
    .catch((err) => {
      if (err.response.status === HTTP_STATUS_CODE.UNPROCESSABLE_ENTITY) {
        showValidationError.value = true;
        invalidEmailMessage.value = err.response.data.message;
      } else if (err.response.status === HTTP_STATUS_CODE.NOT_FOUND) {
        showValidationError.value = true;
        invalidEmailMessage.value = err.response.data.errors[0];
      }
    })
    .finally(() => {
      loading.hide();
    });
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
