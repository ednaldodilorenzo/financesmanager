import { format } from "date-fns";

const months = [
  "Janeiro",
  "Fevereiro",
  "Março",
  "Abril",
  "Maio",
  "Junho",
  "Julho",
  "Agosto",
  "Setembro",
  "Outubro",
  "Novembro",
  "Dezembro",
];

const daysPerMonth = [31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31];

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

function getAproximateMonths(date) {
  let year = date.getUTCFullYear();
  let currentMonth = date.getMonth();

  let result = [];
  for (let i = 1; i < 3; i++) {
    let month = currentMonth + i;
    if (month > 11) {
      month = month - 12;
      year++;
    }
    result.push({
      month: month + 1,
      year: year,
      description: `${months[month]}/${year}`,
    });
  }

  return result;
}

function getDaysListPerMonth(date) {
  if (
    (date.getUTCFullYear() % 4 === 0 && date.getUTCFullYear() % 100 !== 0) ||
    date.getUTCFullYear() % 400 === 0
  ) {
    daysPerMonth[2] = 29;
  }

  return Array.from(
    { length: daysPerMonth[date.getUTCMonth()] },
    (_, i) => 1 + i
  );
}

function getMonthsListUntilDate(date) {
  return months
    .filter((_, index) => index <= date.getUTCMonth())
    .map((value) => value.substring(0, 3));
}

export {
  formatDateUTC,
  getExtenseMonth,
  getDaysListPerMonth,
  getMonthsListUntilDate,
  getAproximateMonths,
};
