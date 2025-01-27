import { formatCurrency } from "@/utils/numbers";

export default {
  beforeMount(el) {
    const input =
      el.tagName === "INPUT" ? el : el.getElementsByTagName("input")[0];

    if (!input) {
      return;
    }

    // Add a data attribute to track if the value has already been formatted
    input.dataset.isFormatted = false;

    input.addEventListener("input", (event) => {
      if (input.dataset.isFormatted === "true") {
        input.dataset.isFormatted = false;
        return;
      }

      let value = input.value;

      // Handle empty input to prevent NaN
      if (value.trim() === "") {
        input.dataset.isFormatted = true; // Set the flag before changing the value
        input.value = ""; // Clear the field
        input.dispatchEvent(new Event("input", { bubbles: true })); // Trigger input event
        return;
      }

      // Check if the value has a negative sign
      //const isNegative = value.startsWith("-");

      // Remove all non-numeric characters except digits
      //value = value.replace(/[^\d]/g, "");

      // Parse the value as a float and divide by 100 for decimal formatting
      //const number = parseFloat(value) / 100;

      // Format the number as Brazilian currency
      //const formattedValue = number.toLocaleString("pt-BR", {
      //  minimumFractionDigits: 2,
      //  maximumFractionDigits: 2,
      //});

      // Add the negative sign back if applicable
      input.dataset.isFormatted = true; // Set the flag before changing the value
      input.value = formatCurrency(input.value);
      //input.value = isNegative ? `-${formattedValue}` : formattedValue;

      // Trigger an input event to update v-model
      input.dispatchEvent(new Event("input", { bubbles: true }));
    });
  },
};
