<template>
  <LoadingScreen :loading="loading" />
  <main class="login-page">
    <div class="row g-0 min-vh-100">
      <div class="col-lg-6 d-none d-lg-flex brand-panel">
        <div class="brand-content">
          <div class="brand-mark mb-4"><i class="bi bi-wallet2"></i></div>
          <span class="eyebrow">FINANCE COCKPIT</span>
          <h1 class="display-5 fw-bold mt-3 mb-4">Suas finanças em uma visão mais clara.</h1>
          <p class="lead text-white-50 mb-5">Planeje, acompanhe e tome decisões melhores com seus dados organizados em
            um só lugar.</p>
          <div class="d-grid gap-3">
            <div class="benefit"><i class="bi bi-check2-circle"></i><span>Acompanhe receitas e despesas</span></div>
            <div class="benefit"><i class="bi bi-check2-circle"></i><span>Planeje seus investimentos</span></div>
            <div class="benefit"><i class="bi bi-check2-circle"></i><span>Visualize sua evolução financeira</span></div>
          </div>
        </div>
        <span class="shape shape-one"></span><span class="shape shape-two"></span>
      </div>

      <div class="col-lg-6 d-flex align-items-center justify-content-center form-panel p-3 p-sm-4 p-xl-5">
        <section class="login-form-wrapper w-100">
          <div class="d-flex align-items-center gap-2 d-lg-none mb-5">
            <span class="mobile-brand-mark"><i class="bi bi-wallet2"></i></span><strong>Finance Cockpit</strong>
          </div>
          <div class="mb-4">
            <span class="small fw-semibold text-primary">BEM-VINDO DE VOLTA</span>
            <h2 class="h3 fw-bold mt-2 mb-2">Entre na sua conta</h2>
            <p class="text-body-secondary mb-0">Informe seus dados para acessar o painel.</p>
          </div>

          <div v-if="showInvalidLoginMessage" class="alert alert-danger d-flex align-items-start gap-3 border-0"
            role="alert" data-test="msg-invalid-login">
            <i class="bi bi-exclamation-circle-fill mt-1"></i>
            <div class="flex-grow-1"><strong class="d-block">Não foi possível entrar</strong><span class="small">Confira
                seu e-mail e senha e tente novamente.</span></div>
            <button type="button" class="btn-close" aria-label="Fechar"
              @click="showInvalidLoginMessage = false"></button>
          </div>

          <form novalidate @submit.prevent="onSubmit">
            <div class="mb-3">
              <label for="inputUsername" class="form-label small fw-semibold">E-mail</label>
              <div class="input-icon-wrapper">
                <i class="bi bi-envelope"></i>
                <input id="inputUsername" v-model.trim="form.email" class="form-control form-control-lg"
                  :class="{ 'is-invalid': v$.email.$error }" type="email" inputmode="email" autocomplete="username"
                  placeholder="voce@exemplo.com" data-test="input-email" @input="showInvalidLoginMessage = false" />
              </div>
              <div v-if="v$.email.$error" class="invalid-message">Informe um e-mail válido.</div>
            </div>

            <div class="mb-3">
              <div class="d-flex align-items-center justify-content-between gap-2">
                <label for="inputPassword" class="form-label small fw-semibold">Senha</label>
                <router-link :to="{ name: ROUTE_NAMES.RECOVER }" class="small text-decoration-none">Esqueceu a
                  senha?</router-link>
              </div>
              <div class="input-icon-wrapper password-wrapper">
                <i class="bi bi-lock"></i>
                <input id="inputPassword" v-model="form.password" class="form-control form-control-lg"
                  :class="{ 'is-invalid': v$.password.$error }" :type="showPassword ? 'text' : 'password'"
                  autocomplete="current-password" placeholder="Informe sua senha" data-test="input-password"
                  @input="showInvalidLoginMessage = false" />
                <button class="password-toggle" type="button"
                  :aria-label="showPassword ? 'Ocultar senha' : 'Mostrar senha'" @click="showPassword = !showPassword">
                  <i :class="['bi', showPassword ? 'bi-eye-slash' : 'bi-eye']"></i>
                </button>
              </div>
              <div v-if="v$.password.$error" class="invalid-message">A senha deve possuir pelo menos 3 caracteres.</div>
            </div>

            <div class="form-check mb-4">
              <input id="rememberUser" v-model="rememberUser" class="form-check-input" type="checkbox" />
              <label class="form-check-label small" for="rememberUser">Lembrar meu e-mail neste dispositivo</label>
            </div>

            <button class="btn btn-primary btn-lg w-100 login-button" data-test="button-login" type="submit"
              :disabled="loading">
              <span v-if="loading" class="spinner-border spinner-border-sm me-2" aria-hidden="true"></span>
              {{ loading ? "Entrando…" : "Entrar" }}
            </button>

            <p class="text-center text-body-secondary mt-4 mb-0">
              Ainda não possui uma conta?
              <router-link :to="{ name: ROUTE_NAMES.SEND_MAIL }" class="fw-semibold text-decoration-none">Cadastre-se
                agora</router-link>
            </p>
          </form>
          <p class="text-center small text-body-tertiary mt-5 mb-0">Seus dados permanecem protegidos.</p>
        </section>
      </div>
    </div>
  </main>
</template>

<script setup>
import { reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { useVuelidate } from "@vuelidate/core";
import { email, minLength, required } from "@vuelidate/validators";
import LoadingScreen from "@/components/loading-screen.vue";
import authService from "./auth.service";
import { ROUTE_NAMES } from "./routes.definition";
import { ROUTE_NAMES as DASHBOARD_ROUTE_NAMES } from "../dashboard/routes.definition";

const REMEMBERED_EMAIL_KEY = "finance-cockpit-email";
const rememberedEmail = localStorage.getItem(REMEMBERED_EMAIL_KEY) || "";
const router = useRouter();
const form = reactive({ email: rememberedEmail, password: "" });
const rememberUser = ref(Boolean(rememberedEmail));
const showPassword = ref(false);
const showInvalidLoginMessage = ref(false);
const loading = ref(false);
const rules = { email: { required, email }, password: { required, minLength: minLength(3) } };
const v$ = useVuelidate(rules, form);

async function onSubmit() {
  showInvalidLoginMessage.value = false;
  if (!await v$.value.$validate()) return;
  loading.value = true;
  try {
    const response = await authService.login(form.email, form.password);
    if (!response) {
      showInvalidLoginMessage.value = true;
      return;
    }
    if (rememberUser.value) localStorage.setItem(REMEMBERED_EMAIL_KEY, form.email);
    else localStorage.removeItem(REMEMBERED_EMAIL_KEY);
    await router.push({ name: DASHBOARD_ROUTE_NAMES.INDEX });
  } catch (error) {
    showInvalidLoginMessage.value = true;
  } finally {
    loading.value = false;
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  background: var(--bs-body-bg)
}

.brand-panel {
  position: relative;
  overflow: hidden;
  align-items: center;
  color: #fff;
  background: linear-gradient(145deg, #3033b8 0%, #4f46e5 50%, #6366f1 100%)
}

.brand-content {
  position: relative;
  z-index: 2;
  width: min(36rem, 80%);
  margin-inline: auto
}

.brand-mark,
.mobile-brand-mark {
  display: grid;
  place-items: center;
  border-radius: 1rem
}

.brand-mark {
  width: 4rem;
  height: 4rem;
  background: rgba(255, 255, 255, .14);
  font-size: 1.8rem
}

.mobile-brand-mark {
  width: 2.5rem;
  height: 2.5rem;
  background: rgba(var(--bs-primary-rgb), .1);
  color: var(--bs-primary)
}

.eyebrow {
  font-size: .72rem;
  font-weight: 700;
  letter-spacing: .18em;
  color: #d9dcff
}

.benefit {
  display: flex;
  align-items: center;
  gap: .75rem;
  color: rgba(255, 255, 255, .82)
}

.benefit i {
  color: #b8f7d4;
  font-size: 1.2rem
}

.shape {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, .06)
}

.shape-one {
  width: 26rem;
  height: 26rem;
  right: -10rem;
  top: -9rem
}

.shape-two {
  width: 20rem;
  height: 20rem;
  left: -8rem;
  bottom: -8rem
}

.form-panel {
  background: #f8fafc
}

.login-form-wrapper {
  max-width: 28rem
}

.alert {
  border-radius: .75rem;
  background: rgba(var(--bs-danger-rgb), .1);
  color: var(--bs-danger-text-emphasis)
}

.input-icon-wrapper {
  position: relative
}

.input-icon-wrapper>i {
  position: absolute;
  left: 1rem;
  top: 50%;
  z-index: 2;
  transform: translateY(-50%);
  color: var(--bs-secondary-color)
}

.input-icon-wrapper .form-control {
  padding-left: 2.75rem;
  border-color: var(--bs-border-color);
  background: #fff
}

.password-wrapper .form-control {
  padding-right: 3rem
}

.password-toggle {
  position: absolute;
  right: .65rem;
  top: 50%;
  z-index: 4;
  transform: translateY(-50%);
  width: 2.25rem;
  height: 2.25rem;
  border: 0;
  border-radius: 50%;
  background: transparent;
  color: var(--bs-secondary-color)
}

.password-toggle:hover {
  background: var(--bs-tertiary-bg);
  color: var(--bs-body-color)
}

.invalid-message {
  margin-top: .35rem;
  color: var(--bs-danger);
  font-size: .78rem
}

.login-button {
  min-height: 3.15rem;
  border-radius: .65rem;
  box-shadow: 0 .5rem 1.25rem rgba(var(--bs-primary-rgb), .2)
}

@media(max-width:991.98px) {
  .form-panel {
    background: linear-gradient(180deg, rgba(var(--bs-primary-rgb), .045), #fff 35%)
  }
}
</style>