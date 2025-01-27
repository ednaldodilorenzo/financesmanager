function parseCurrencyToNumber(formattedValue) {
  if (!formattedValue) return 0;
  const sanitizedValue = formattedValue.replace(/[^\d,-]/g, "");
  const numericValue = sanitizedValue.replace(",", "");
  return parseInt(numericValue) || 0;
}

const formatCurrency = (value) => {
  // Check if the value has a negative sign
  const isNegative = value.startsWith("-");

  // Remove all non-numeric characters except digits
  const digits = value.replace(/[^\d]/g, "");

  // Parse the value as a float and divide by 100 for decimal formatting
  const number = digits ? parseFloat(digits) / 100 : 0;

  // Format the number as Brazilian currency
  const formattedValue = number.toLocaleString("pt-BR", {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  });

  return isNegative ? `-${formattedValue}` : formattedValue;
};

export { parseCurrencyToNumber, formatCurrency };
