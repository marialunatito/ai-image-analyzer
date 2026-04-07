export default function StatusMessage({ variant = "info", children }) {
  if (!children) {
    return null;
  }

  return (
    <p className={`status-message status-${variant}`} role="status" aria-live="polite">
      {children}
    </p>
  );
}
