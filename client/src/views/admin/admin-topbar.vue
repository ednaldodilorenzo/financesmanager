<template>
  <ul class="nav navbar justify-content-end bg-body-tertiary shadow px-3">
    <li class="nav-item">
      <div class="dropdown">
        <button
          class="btn btn-icon btn-topbar btn-ghost-secondary rounded-circle dropdown-toggle"
          :class="{ show: userMenuToogle }"
          type="button"
          id="dropdownMenuButton"
          data-bs-toggle="dropdown"
          aria-expanded="false"
          @click="toggleUserMenu"
        >
          <i class="bi bi-person" style="font-size: 1.5rem"></i>
        </button>
        <ul
          class="dropdown-menu dropdown-menu-end"
          :class="{ show: userMenuToogle }"
          aria-labelledby="dropdownMenuButton"
          v-bind="{ attribute: condition ? value : null }"
          :data-bs-popper="userMenuToogle ? 'static' : null"
        >
          <li>
            <a class="dropdown-item" @click="onChangePasswordClick" href="#"
              >Alterar Senha</a
            >
          </li>
        </ul>
      </div>
    </li>
    <li class="nav-item">
      <button
        type="button"
        class="btn btn-icon btn-topbar btn-ghost-secondary rounded-circle"
        data-toggle="fullscreen"
        @click="logout()"
      >
        <i class="bi bi-box-arrow-right" style="font-size: 1.5rem"></i>
      </button>
    </li>
  </ul>
</template>
<script setup>
import { ref } from "vue";
import authService from "@/views/login/auth.service";
import { ROUTE_NAMES as AUTH_ROUTES } from "@/views/login/routes.definition";
import { useRouter } from "vue-router";
import { useModalScreen } from "@/components/modal/use-modal-screen";
import ChangePassword from "../login/change-password.vue";

const emit = defineEmits(["collapse-click"]);
const router = useRouter();
const userMenuToogle = ref(false);
const modal = useModalScreen(ChangePassword);

const logout = () => {
  authService.logout().then(() => {
    router.push({ name: AUTH_ROUTES.INDEX });
  });
};

const onChangePasswordClick = async () => {
  toggleUserMenu();
  const confirm = await modal.show();  
};

const toggleUserMenu = () => {
  userMenuToogle.value = !userMenuToogle.value;
};
</script>
