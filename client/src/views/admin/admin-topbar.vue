<template>
  <nav class="navbar bg-body-tertiary shadow">
    <div class="container-fluid">
      <button
        type="button"
        class="btn btn-icon btn-topbar btn-ghost-secondary rounded-circle"
        @click="toggleMenu()"
        data-toggle="fullscreen"
      >
        <i class="bx bx-menu bx-sm"></i>
      </button>
      <form class="d-flex" role="search">
        <input
          class="form-control me-2"
          type="search"
          placeholder="Search"
          aria-label="Search"
        />
        <button class="btn btn-outline-success" type="submit">Search</button>
      </form>
      <div>
        <button
          type="button"
          class="btn btn-icon btn-topbar btn-ghost-secondary rounded-circle"
          data-toggle="fullscreen"
        >
          <i class="bx bx-fullscreen fs-4"></i>
        </button>
        <button
          type="button"
          class="btn btn-icon btn-topbar btn-ghost-secondary rounded-circle"
          data-toggle="fullscreen"
        >
          <i class="bx bx-bell fs-4"></i>
        </button>
        <button
          type="button"
          class="btn btn-icon btn-topbar btn-ghost-secondary rounded-circle"
          data-toggle="fullscreen"
          @click="logout()"
        >
          <i
            class="bi bi-box-arrow-right"
            style="font-size: 1.5rem"
          ></i>
        </button>
      </div>
    </div>
  </nav>
</template>
<script>
import authService from "@/views/login/auth.service";
import { ROUTE_NAMES as AUTH_ROUTES } from "@/views/login/routes.definition";

export default {
  name: "admin-topbar",
  emits: ["collapse-click"],
  methods: {
    logout() {
      authService.logout().then(() => {
        this.$router.push({ name: AUTH_ROUTES.INDEX });
      });
    },
    toggleMenu() {
      this.$emit("collapse-click");
    },
  },
};
</script>
