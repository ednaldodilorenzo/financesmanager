<template>redirecionando...</template>
<script setup>
import { useRoute } from "vue-router";
import { useRouter } from "vue-router";
import authService from "./auth.service";
import { ROUTE_NAMES as DASHBOARD_ROUTE_NAMES } from "@/views/dashboard/routes.definition";

const route = useRoute();
const router = useRouter();

authService.oauthLogin(route.params.userId, route.params.name).then(() => {
  router.push({ name: DASHBOARD_ROUTE_NAMES.INDEX });
});
</script>
