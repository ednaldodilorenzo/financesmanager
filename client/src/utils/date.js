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

function getDaysListPerMonth(date) {
  const daysPerMonth = [31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31];

  if (((date.getUTCFullYear() % 4 === 0) && (date.getUTCFullYear() % 100 !== 0)) || (date.getUTCFullYear() % 400 === 0)) {
    daysPerMonth[2] = 29  
  }

  return Array.from({ length: daysPerMonth[date.getUTCMonth()] }, (_, i) => 1 + i);
}

export { formatDateUTC, getExtenseMonth, getDaysListPerMonth };
