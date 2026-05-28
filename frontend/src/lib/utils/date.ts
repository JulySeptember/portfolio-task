export function formatDate(date: string | Date, locale: string = "ja-JP") {
  const value = typeof date === "string" ? new Date(date) : date;

  return new Intl.DateTimeFormat(locale, {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  }).format(value);
}

export function formatDateOnly(date: string | Date, locale: string = "ja-JP") {
  const value = typeof date === "string" ? new Date(date) : date;

  return new Intl.DateTimeFormat(locale, {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
  }).format(value);
}

export function isPast(date: string | Date) {
  const value = typeof date === "string" ? new Date(date) : date;

  return value.getTime() < Date.now();
}
