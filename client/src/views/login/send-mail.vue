<template>
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
  {{ message }}
  <div class="mt-4 text-center">
    <p class="mb-0">
      <label class="form-label">Informe o email</label>
      <input
        type="text"
        v-model="email"
        class="form-control my-3"
        placeholder="Email"
      />
      <button @click="sendEmail" type="button" class="btn btn-primary">
        Enviar
      </button>
    </p>
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
const loading = useLoadingScreen();
const toast = useToast();
const route = useRoute();

const validationMessage = ref(null);
const showValidationError = ref(false);

const sendEmail = () => {
  loading.show();
  const method =
    route.name === ROUTE_NAMES.RECOVER
      ? authService.startRecoverProcess(email.value)
      : authService.startRegistration(email.value);
  method
    .then(() => {
      toast.success("Email enviado com sucesso!");
    })
    .catch((err) => {
      if (err.response.status === HTTP_STATUS_CODE.UNPROCESSABLE_ENTITY) {
        showValidationError.value = true;
        validationMessage.value = err.response.data.message;
      }
    })
    .finally(() => {
      loading.hide();
    });
};
</script>
