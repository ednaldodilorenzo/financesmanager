<template>
  <section class="import-page py-3">
    <header class="mb-3">
      <h1 class="h4 fw-bold mb-1">Importar transações</h1>
      <p class="small text-body-secondary mb-0">Envie um extrato, revise os lançamentos e escolha o que deseja importar
      </p>
    </header>

    <div class="row g-4 align-items-start">
      <div class="col-xl-4">
        <form class="card border-0 shadow-sm import-config sticky-xl-top" @submit.prevent="submitForm">
          <div class="card-header bg-white border-bottom p-3">
            <div class="d-flex align-items-center gap-3">
              <span class="step-icon">1</span>
              <div>
                <h2 class="h6 fw-bold mb-1">Selecionar arquivo</h2>
                <p class="small text-body-secondary mb-0">Configure a origem da importação</p>
              </div>
            </div>
          </div>
          <div class="card-body p-3">
            <label for="fileType" class="form-label small fw-semibold">Formato do arquivo</label>
            <select id="fileType" v-model="fileType" class="form-select mb-3">
              <option value="BBCA">Conta Corrente BB</option>
              <option value="C6CC">Cartão de Crédito C6</option>
              <option value="CUAL">Arquivo customizado</option>
            </select>

            <label v-if="requiresAccount" for="accountId" class="form-label small fw-semibold">Conta de destino</label>
            <select v-if="requiresAccount" id="accountId" v-model="currentAccount" class="form-select mb-3">
              <option :value="null" disabled>Selecione uma conta</option>
              <option v-for="account in accounts" :key="account.id" :value="account.id">{{ account.name }}</option>
            </select>

            <template v-if="fileType === 'C6CC'">
              <label for="paymentDate" class="form-label small fw-semibold">Mês da fatura</label>
              <input id="paymentDate" v-model="paymentDate" type="month" class="form-control mb-3" />
            </template>

            <label for="transactionFile" class="upload-zone" :class="{ 'has-file': formData.file }">
              <input id="transactionFile" class="visually-hidden" type="file" accept=".csv,.ofx,.txt"
                @change="handleFileUpload" />
              <span class="upload-icon"><i
                  :class="formData.file ? 'bi bi-file-earmark-check' : 'bi bi-cloud-arrow-up'"></i></span>
              <strong class="d-block mt-2">{{ formData.file?.name || 'Escolha o arquivo' }}</strong>
              <span class="small text-body-secondary">CSV, OFX ou TXT</span>
            </label>
          </div>
          <div class="card-footer bg-white border-top p-3">
            <button type="submit" class="btn btn-primary w-100" :disabled="!canPrepare">
              <i class="bi bi-search me-1"></i>Preparar importação
            </button>
          </div>
        </form>
      </div>

      <div class="col-xl-8">
        <div class="card border-0 shadow-sm review-card">
          <div class="card-header bg-white border-bottom p-3">
            <div class="d-flex flex-wrap align-items-center justify-content-between gap-3">
              <div class="d-flex align-items-center gap-3">
                <span class="step-icon">2</span>
                <div>
                  <h2 class="h6 fw-bold mb-1">Revisar lançamentos</h2>
                  <p class="small text-body-secondary mb-0">Confira as informações antes de importar</p>
                </div>
              </div>
              <div v-if="collection.length" class="d-flex flex-wrap gap-2">
                <span class="badge rounded-pill text-bg-light">{{ collection.length }} pendentes</span>
                <span v-if="selectedCount" class="badge rounded-pill text-bg-primary">{{ selectedCount }}
                  selecionados</span>
                <span v-if="duplicatedCount" class="badge rounded-pill text-bg-danger">{{ duplicatedCount }}
                  duplicados</span>
              </div>
            </div>
          </div>

          <div v-if="collection.length" class="batch-toolbar border-bottom p-3">
            <div class="d-flex flex-wrap align-items-center gap-2">
              <button class="btn btn-primary btn-sm" type="button" @click="importNonDuplicated">
                <i class="bi bi-box-arrow-in-down me-1"></i>Importar não duplicados
              </button>
              <button class="btn btn-outline-primary btn-sm" type="button" :disabled="!selectedCount"
                @click="importSelected">
                Importar selecionados
              </button>
              <button class="btn btn-outline-danger btn-sm ms-sm-auto" type="button" :disabled="!duplicatedCount"
                @click="discardDuplicated">
                <i class="bi bi-trash3 me-1"></i>Descartar duplicados
              </button>
            </div>
          </div>

          <div v-if="collection.length" class="transaction-review-list">
            <ValidateEach v-for="item in collection" :key="item._key" :state="item" :rules="rules">
              <template #default="{ v }">
                <article class="import-row p-3" :class="{ 'is-duplicated': item.duplicated }">
                  <div class="d-flex align-items-start gap-3">
                    <input v-model="item.checked" class="form-check-input mt-2" type="checkbox"
                      :aria-label="`Selecionar ${item.description}`" />
                    <div class="flex-grow-1 min-w-0">
                      <div class="d-flex flex-wrap align-items-center justify-content-between gap-2 mb-3">
                        <div class="d-flex align-items-center gap-2">
                          <span :class="['transaction-type-icon', item.duplicated ? 'is-duplicate' : '']"><i
                              class="bi bi-receipt"></i></span>
                          <div>
                            <span class="small text-body-secondary d-block">{{ formatDisplayDate(item.paymentDate)
                              }}</span>
                            <span v-if="item.duplicated" class="badge text-bg-danger">Possível duplicidade</span>
                          </div>
                        </div>
                        <strong :class="amountClass(item.formatted_value)">{{ item.formatted_value }}</strong>
                      </div>

                      <div class="row g-3">
                        <div class="col-lg-6">
                          <label class="form-label small fw-semibold">Descrição</label>
                          <input v-model="v.description.$model" type="text" class="form-control"
                            :class="{ 'is-invalid': v.description.$errors.length }" />
                          <div class="invalid-feedback">Informe a descrição.</div>
                        </div>
                        <div class="col-lg-6">
                          <label class="form-label small fw-semibold">Categoria</label>
                          <select v-model="v.category.$model" class="form-select"
                            :class="{ 'is-invalid': v.category.$errors.length }">
                            <option value="" disabled>Selecione uma categoria</option>
                            <option v-for="category in categories" :key="category.id" :value="category.id">{{
                              category.name }}</option>
                          </select>
                          <div class="invalid-feedback">Selecione uma categoria.</div>
                        </div>
                        <div class="col-lg-6">
                          <label class="form-label small fw-semibold">Anotação</label>
                          <input v-model="item.detail" type="text" class="form-control" placeholder="Opcional" />
                        </div>
                        <div class="col-lg-6">
                          <label class="form-label small fw-semibold">Tags</label>
                          <BootstrapSelectTag v-model="item.tags" :options="searchTags" :id="`tags-${item._key}`"
                            name="tags" />
                        </div>
                        <div class="col-sm-6 col-lg-4">
                          <label class="form-label small fw-semibold">Data</label>
                          <input v-model="v.formatted_date.$model" type="date" class="form-control" disabled />
                        </div>
                        <div class="col-sm-6 col-lg-4">
                          <label class="form-label small fw-semibold">Valor</label>
                          <div class="input-group">
                            <span class="input-group-text">R$</span>
                            <input v-model="v.formatted_value.$model" v-currency type="text"
                              class="form-control text-end"
                              :class="{ 'is-invalid': v.formatted_value.$errors.length }" />
                          </div>
                        </div>
                        <div class="col-lg-4 d-flex align-items-end justify-content-lg-end gap-2">
                          <button class="btn btn-success" type="button" @click="importTransaction(v, item)">
                            <i class="bi bi-check-lg me-1"></i>Importar
                          </button>
                          <button class="btn btn-outline-danger" type="button" title="Descartar"
                            @click="removeItem(item)">
                            <i class="bi bi-x-lg"></i>
                          </button>
                        </div>
                      </div>
                    </div>
                  </div>
                </article>
              </template>
            </ValidateEach>
          </div>

          <div v-else class="text-center py-5 px-3">
            <span class="empty-icon"><i class="bi bi-file-earmark-arrow-up"></i></span>
            <h3 class="h6 mt-3 mb-1">Nenhum lançamento para revisar</h3>
            <p class="small text-body-secondary mb-0">Selecione um arquivo e prepare a importação.</p>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed, ref } from "vue";
import { useVuelidate } from "@vuelidate/core";
import { required } from "@vuelidate/validators";
import { ValidateEach } from "@vuelidate/components";
import transactionService from "./transaction.service";
import categoryService from "../category/category.service";
import accountService from "../account/account.service";
import tagService from "@/views/tag/tag.service";
import BootstrapSelectTag from "@/components/bootstrap-select-tag.vue";
import { useLoadingScreen } from "@/components/loading/useLoadingScreen";
import { formatDateUTC } from "@/utils/date";
import { parseCurrencyToNumber } from "@/utils/numbers";
import { useToast } from "vue-toastification";
import { currencyBRL } from "@/components/filters/currency.filter";

const loading = useLoadingScreen();
const toast = useToast();
const currentAccount = ref(null);
const categories = ref([]);
const accounts = ref([]);
const fileType = ref("BBCA");
const paymentDate = ref("");
const formData = ref({ file: null });
const state = ref({ collection: [] });
const v$ = useVuelidate();

const collection = computed(() => state.value.collection);
const selectedCount = computed(() => collection.value.filter(({ checked }) => checked).length);
const duplicatedCount = computed(() => collection.value.filter(({ duplicated }) => duplicated).length);
const requiresAccount = computed(() => ["BBCA", "C6CC"].includes(fileType.value));
const canPrepare = computed(() => Boolean(formData.value.file && (!requiresAccount.value || currentAccount.value) && (fileType.value !== "C6CC" || paymentDate.value)));

const rules = {
  formatted_date: { required },
  description: { required },
  category: { required },
  formatted_value: { required },
};

async function getDependencies() {
  loading.show();
  try {
    const [categoryResponse, accountResponse] = await Promise.all([
      categoryService.findAll({ paginate: false }),
      accountService.findAll({ paginate: false }),
    ]);
    categories.value = categoryResponse.data;
    accounts.value = accountResponse.data;
  } finally { loading.hide(); }
}

const searchTags = (filter) => tagService.findAll({ filter }).then((response) => response.items.map(({ tag }) => tag));
const formatDisplayDate = (date) => formatDateUTC(date, "dd/MM/yyyy");
const amountClass = (value) => parseCurrencyToNumber(value) < 0 ? "text-danger" : "text-primary";

function removeItem(clicked) {
  state.value.collection = collection.value.filter((item) => item !== clicked);
}

async function importNonDuplicated() {
  if (!await v$.value.$validate()) return;
  sendBatchData(collection.value.filter(({ duplicated }) => !duplicated));
}

async function importSelected() {
  if (!await v$.value.$validate()) return;
  sendBatchData(collection.value.filter(({ checked }) => checked));
}

async function importTransaction(validator, item) {
  if (!await validator.$validate()) return;
  loading.show();
  try {
    await transactionService.create(toPayload(item));
    removeItem(item);
    toast.success("Transação criada com sucesso!", { position: "top-center" });
  } finally { loading.hide(); }
}

function discardDuplicated() {
  state.value.collection = collection.value.filter(({ duplicated }) => !duplicated);
}

function toPayload(item) {
  const payment = new Date(item.paymentDate);
  return {
    categoryId: item.category,
    accountId: currentAccount.value,
    detail: item.detail,
    description: item.description,
    tags: (item.tags || []).map((tag) => ({ tag })),
    value: parseCurrencyToNumber(item.formatted_value),
    paymentDate: payment.toISOString(),
    paymentMonth: payment.getUTCMonth() + 1,
    paymentYear: payment.getUTCFullYear(),
    transactionDate: new Date(item.transactionDate).toISOString(),
  };
}

async function sendBatchData(items) {
  if (!items.length) return;
  loading.show();
  try {
    await transactionService.sendBatchImport(items.map(toPayload));
    state.value.collection = collection.value.filter((item) => !items.includes(item));
    toast.success(`${items.length} transações importadas com sucesso!`, { position: "top-center" });
  } finally { loading.hide(); }
}

function handleFileUpload(event) {
  formData.value.file = event.target.files?.[0] || null;
}

async function submitForm() {
  if (!canPrepare.value) {
    toast.warning("Preencha os dados da importação e selecione um arquivo.", { position: "top-center" });
    return;
  }
  const payload = new FormData();
  payload.append("file", formData.value.file);
  payload.append("fileType", fileType.value);
  if (requiresAccount.value) payload.append("accountId", currentAccount.value);
  if (fileType.value === "C6CC") {
    const [year, month] = paymentDate.value.split("-");
    payload.append("paymentMonth", month);
    payload.append("paymentYear", year);
  }
  loading.show();
  try {
    const response = await transactionService.prepareForImport(payload);
    state.value.collection = response.data.map((item, index) => ({
      ...item,
      _key: `${item.transactionDate}-${item.description}-${index}`,
      formatted_date: formatDateUTC(item.paymentDate, "yyyy-MM-dd"),
      formatted_value: currencyBRL(item.value),
      category: item.categoryId || "",
      checked: false,
      tags: [],
    }));
  } finally { loading.hide(); }
}

getDependencies();
</script>

<style scoped>
.import-config,
.review-card {
  border-radius: .85rem;
}

.import-config {
  top: 1rem;
}

.import-config .card-header,
.review-card .card-header {
  border-radius: .85rem .85rem 0 0;
}

.step-icon {
  width: 2rem;
  height: 2rem;
  display: grid;
  place-items: center;
  flex: 0 0 auto;
  border-radius: 50%;
  background: rgba(var(--bs-primary-rgb), .12);
  color: var(--bs-primary);
  font-weight: 700;
}

.upload-zone {
  display: block;
  padding: 1.5rem 1rem;
  border: 2px dashed var(--bs-border-color);
  border-radius: .75rem;
  text-align: center;
  cursor: pointer;
  transition: border-color .15s ease, background-color .15s ease;
}

.upload-zone:hover,
.upload-zone.has-file {
  border-color: var(--bs-primary);
  background: rgba(var(--bs-primary-rgb), .035);
}

.upload-icon,
.empty-icon {
  display: inline-grid;
  place-items: center;
  background: rgba(var(--bs-primary-rgb), .1);
  color: var(--bs-primary);
}

.upload-icon {
  width: 2.75rem;
  height: 2.75rem;
  border-radius: .8rem;
  font-size: 1.25rem;
}

.empty-icon {
  width: 3.25rem;
  height: 3.25rem;
  border-radius: 1rem;
  font-size: 1.4rem;
}

.batch-toolbar {
  background: var(--bs-tertiary-bg);
}

.import-row {
  border-bottom: 1px solid var(--bs-border-color-translucent);
}

.import-row:last-child {
  border-bottom: 0;
}

.import-row.is-duplicated {
  border-left: 3px solid var(--bs-danger);
  background: rgba(var(--bs-danger-rgb), .025);
}

.transaction-type-icon {
  width: 2.25rem;
  height: 2.25rem;
  display: grid;
  place-items: center;
  border-radius: .65rem;
  background: rgba(var(--bs-primary-rgb), .1);
  color: var(--bs-primary);
}

.transaction-type-icon.is-duplicate {
  background: rgba(var(--bs-danger-rgb), .1);
  color: var(--bs-danger);
}

.min-w-0 {
  min-width: 0;
}

@media (max-width: 1199.98px) {
  .import-config {
    position: static !important;
  }
}
</style>