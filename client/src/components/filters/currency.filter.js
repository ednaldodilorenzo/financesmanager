export const currencyBRL = (value) => {
  return Intl.NumberFormat("pt-br", {
    style: "currency",
    currency: "BRL",
  }).format(value / 100);
};

export const percentageBRL = (value) => {
  if (!value) return "";
  return Intl.NumberFormat("pt-br", {
    style: "percent",
    currency: "BRL",
  }).format(value);
};
