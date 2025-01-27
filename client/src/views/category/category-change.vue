<template>
  <bootstrap-modal
    @close="onCloseModal()"
    @save="onSaveModal()"
    :title="`${form.id ? 'Alterar' : 'Nova'} Categoria`"
    :visible="open"
  >
    <form
      :class="{
        'was-validated': v$.$dirty ? true : false,
      }"
      @submit.stop.prevent="onSubmit"
      id="frmCategory"
      class="row g-3 mb-3"
      autocomplete="off"
      novalidate
    >
      <div class="col-md-6">
        <bootstrap-input
          required-message="Por favor preencha o nome da categoria"
          :required="true"
          v-model="form.name"
          id="iptNome"
          label="Nome"
          name="nome"
          class="form-control"
        />
      </div>
      <div class="col-md-6">
        <bootstrap-select
          required-message="Por favor preencha o tipo da
        categoria"
          :required="true"
          v-model="form.type"
          id="slcTipo"
          :options="[
            { id: '', description: 'Selecione...' },
            { id: 'E', description: 'Receita' },
            { id: 'X', description: 'Despesa' },
          ]"
          :keyField="'id'"
          :valueField="'description'"
          label="Tipo"
          class="form-control"
        />
      </div>
    </form>
  </bootstrap-modal>
</template>
<script>
import BootstrapModal from "@/components/bootstrap-modal.vue";
import BootstrapInput from "@/components/bootstrap-input.vue";
import BootstrapSelect from "@/components/bootstrap-select.vue";
import { useVuelidate } from "@vuelidate/core";
import categoryService from "./category.service";
import { required } from "@vuelidate/validators";
import { useToast } from "vue-toastification";

export default {
  setup() {
    const toast = useToast();
    return { v$: useVuelidate(), toast };
  },
  components: { BootstrapModal, BootstrapInput, BootstrapSelect },
  emits: ["close", "saved"],
  props: {
    visible: {
      type: Boolean,
      default: false,
    },
    itemSelected: {
      type: Object,
      default: {
        name: "",
        type: { id: "" },
      },
    },
  },
  data() {
    return {
      loading: false,
      form: this.itemSelected,
      open: this.visible,
    };
  },
  methods: {
    onSaveModal() {
      this.v$.$validate();

      if (this.v$.$error) {
        return;
      }

      this.loading = true;
      const payload = { ...this.form, type: this.form.type.id };

      const method = this.form.id
        ? categoryService.modify(this.form.id, payload)
        : categoryService.create(payload);

      method
        .then(() => {
          this.toast.success(
            `Categoria ${this.form.id ? "atualizada" : "criada"} com sucesso!`,
            {
              position: "top-center",
              timeout: 5000,
              closeOnClick: true,
              pauseOnFocusLoss: true,
              pauseOnHover: true,
              draggable: true,
              draggablePercent: 0.6,
              showCloseButtonOnHover: false,
              hideProgressBar: true,
              closeButton: "button",
              icon: true,
              rtl: false,
            }
          );
          this.v$.$reset();
          this.$emit("saved");
        })
        .catch((e) => {
          console.error(e);
        });
    },
    onCloseModal() {
      this.v$.$reset();
      this.$emit("close");
    },
  },
  validations() {
    return {
      form: {
        name: {
          required,
        },
        type: {
          required,
        },
      },
    };
  },
  watch: {
    visible: function (newVal, oldVal) {
      this.open = newVal;
    },
    itemSelected: function (newVal, oldVal) {
      this.form = {
        ...newVal,
        type: {
          id: newVal.type,
          description: newVal.type === "X" ? "Despesa" : "Receita",
        },
      };
    },
  },
};
</script>
