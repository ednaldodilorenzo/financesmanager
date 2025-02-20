<template>
  <div style="display: inline-flex">
    <a
      href="javascript:void(0)"
      @click="previousMonth"
      class="btn btn-sm icon-link icon-link-hover"
      ><i class="bi bi-chevron-left"></i
    ></a>
    <div class="p-2 date-picker">
      {{ props.onlyYears ? "" : months[currentDate.getMonth()] }}
      {{ currentDate.getFullYear() }}
    </div>
    <a
      href="javascript:void(0)"
      @click="nextMonth"
      class="btn btn-sm icon-link icon-link-hover"
      ><i class="bi bi-chevron-right"></i
    ></a>
  </div>
</template>
<script setup>
import { ref, computed, watch } from "vue";

const months = {
  0: "Janeiro",
  1: "Fevereiro",
  2: "Março",
  3: "Abril",
  4: "Maio",
  5: "Junho",
  6: "Julho",
  7: "Agosto",
  8: "Setembro",
  9: "Outubro",
  10: "Novembro",
  11: "Dezembro",
};

const props = defineProps({
  modelValue: {
    type: Date,
    default: () => new Date(),
  },
  onlyYears: {
    type: Boolean,
    default: false,
  },
});

const currentDate = ref(new Date());

const emit = defineEmits(["update:modelValue", "date-change"]);

watch(
  () => props.modelValue,
  (newValue) => {
    currentDate.value = newValue;
  }
);

const days = computed(() => {
  const daysInMonth = new Date(
    currentYear.value,
    currentMonth.value + 1,
    0
  ).getDate();
  const daysArray = [];
  for (let day = 1; day <= daysInMonth; day++) {
    daysArray.push({
      day,
      date: `${currentYear.value}-${currentMonth.value + 1}-${day}`,
    });
  }
  return daysArray;
});

const currentMonth = computed(() => {
  return currentDate.value.getMonth();
});

const currentYear = computed(() => {
  console.log(currentDate.value.getYear());
  return currentDate.value.getYear();
});

const currentMonthName = computed(() =>
  today.toLocaleString("default", { month: "long" })
);

const nextMonth = () => {
  const newDate = new Date(currentDate.value);
  if (props.onlyYears) {
    newDate.setFullYear(newDate.getFullYear() + 1);
    currentDate.value = newDate;
  } else {
    if (newDate.getMonth() < 11) {
      newDate.setMonth(newDate.getMonth() + 1, 2);
      currentDate.value = newDate;
    } else {
      newDate.setFullYear(newDate.getFullYear() + 1);
      newDate.setMonth(0);
      currentDate.value = newDate;
    }
  }

  emit("date-change", currentDate.value);
  emit("update:modelValue", currentDate.value);
};

const previousMonth = () => {
  const newDate = new Date(currentDate.value);
  if (props.onlyYears) {
    newDate.setFullYear(newDate.getFullYear() - 1);
    currentDate.value = newDate;
  } else {
    if (newDate.getMonth() > 0) {
      newDate.setMonth(newDate.getMonth() - 1);
      currentDate.value = newDate;
    } else {
      newDate.setMonth(11);
      newDate.setFullYear(newDate.getFullYear() - 1);
      currentDate.value = newDate;
    }
  }

  emit("date-change", currentDate.value);
  emit("update:modelValue", currentDate.value);
};
</script>
<style scoped>
.date-picker {
  border: solid black 1px;
  border-radius: 50px;
  text-align: center;
  width: 10rem;
}
</style>
