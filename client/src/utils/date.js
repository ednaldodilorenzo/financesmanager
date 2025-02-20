import { format } from "date-fns";

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

function formatDateUTC(dateString, pattern) {
  const date = new Date(dateString);
  return format(
    new Date(date.getUTCFullYear(), date.getUTCMonth(), date.getUTCDate()),
    pattern
  ); // Ensures the correct day
}

function getExtenseMonth(date) {
  return months[date.getUTCMonth()];
}

export { formatDateUTC, getExtenseMonth };
