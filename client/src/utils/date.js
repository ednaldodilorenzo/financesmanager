function formatDate(date, locale = "pt-br", options = {}) {
  const formatter = new Intl.DateTimeFormat(locale, {
    year: "numeric",
    month: "long",
    day: "numeric",
    ...options,
  });
  return formatter.format(new Date(date));
}

export { formatDate };
