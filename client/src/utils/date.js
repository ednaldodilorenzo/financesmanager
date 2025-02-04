import { format } from "date-fns";

function formatDateUTC(dateString, pattern) {
  const date = new Date(dateString);
  return format(
    new Date(date.getUTCFullYear(), date.getUTCMonth(), date.getUTCDate()),
    pattern
  ); // Ensures the correct day
}

export { formatDateUTC };
