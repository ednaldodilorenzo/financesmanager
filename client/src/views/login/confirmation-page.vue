<template>
    {{ message }}
    <div class="mt-4 text-center">
      <p class="mb-0">
        <router-link to="/login" class="fw-medium text-primary">
          Fazer Login</router-link
        >
      </p>
    </div>
</template>
<script setup>
import authService from './auth.service';
import { ref } from 'vue';
import { useRoute } from 'vue-router';
import { useLoadingScreen } from '@/components/loading/useLoadingScreen';

const route = useRoute();
const loading = useLoadingScreen();

const message = ref(null);

const sendConfirmation = () => {
    loading.show();
    const token = route.params.token;
    authService.confirmAccount(token).then(() => {
        message.value = "Parabéns! Sua inscrição foi confirmada com sucesso.";
    }).catch(err => {
        message.value = "Falha na verificação da inscrição!"
    }).finally(() => {
        loading.hide();
    });
}

sendConfirmation();
</script>