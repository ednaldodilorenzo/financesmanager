<template>
  <bootstrap-modal-screen
    :onClose="onCancelModal"
    :onConfirm="onSubmit"
    size="md"
    title="Alterar Senha"
  >
    <form
      :class="{
        'was-validated': v.$dirty ? true : false,
      }"
      id="frmChangePassword"
      class="row g-3 mb-3"
      autocomplete="off"
    >
      <div class="col-md-12">
        <bootstrap-input
          type="password"
          required
          requiredMessage="Informe a senha atual"
          id="iptPassword"
          v-model="form.password"
          label="Senha Atual"
          name="currentPassword"
        />
      </div>
      <div class="col-md-12">
        <bootstrap-input
          type="password"
          required
          requiredMessage="Informe a nova senha"
          v-model="form.newPassword"
          id="iptNewPassword"
          label="Nova Senha"
          name="newPassword"
        />
      </div>
      <div class="col-md-12">
        <bootstrap-input
          id="iptConfirmPassword"
          required
          requiredMessage="Confirme a nova senha"
          v-model="form.confirmNewPassword"
          label="Confirme Nova Senha"
          name="confirmPassword"
          type="password"
        />
      </div>
    </form>
  </bootstrap-modal-screen>
</template>
<script setup>
import { ref } from "vue";
import BootstrapInput from "@/components/bootstrap-input.vue";
import BootstrapModalScreen from "@/components/bootstrap-modal-screen.vue";
import { required } from "@vuelidate/validators";
import { useToast } from "vue-toastification";
import useVuelidate from "@vuelidate/core";
import authService from "./auth.service";
import { useLoadingScreen } from "@/components/loading/useLoadingScreen";

const props = defineProps({
  onSaveModal: Function,
  onCancelModal: Function,
  item: {
    type: Object,
    required: false,
    default: () => ({
      password: "",
      newPassword: "",
      confirmNewPassword: "",
    }),
  },
});

const rules = {
  password: { required },
  newPassword: { required },
  confirmNewPassword: { required },
};

const form = ref({});
form.value = props.item;
const toast = useToast();
const v = useVuelidate(rules, form.value);
const loading = useLoadingScreen();

const onSubmit = () => {
  v.value.$validate();

  if (v.value.$error) {
    return;
  }

  loading.show();

  authService
    .changePassword(form.value)
    .then(() => {
      toast.success("Senha alterada com sucesso!");
    })
    .finally(() => {
      loading.hide();
    });
};
</script>
