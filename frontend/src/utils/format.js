export function toPercent(value) {
  const parsed = Number(value);
  if (Number.isNaN(parsed)) {
    return "N/A";
  }

  return `${(parsed * 100).toFixed(1)}%`;
}
