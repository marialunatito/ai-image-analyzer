const DEFAULT_TIMEOUT_MS = 20000;

function buildApiUrl(path) {
  const base = import.meta.env.VITE_API_BASE_URL || "";
  return `${base}${path}`;
}

async function parseApiError(response) {
  const fallbackMessage = `Error del servidor (${response.status}).`;
  try {
    const data = await response.json();
    return data?.message || fallbackMessage;
  } catch {
    return fallbackMessage;
  }
}

export async function analyzeImage(file) {
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), DEFAULT_TIMEOUT_MS);

  try {
    const formData = new FormData();
    formData.append("image", file);

    const response = await fetch(buildApiUrl("/api/analyze"), {
      method: "POST",
      body: formData,
      signal: controller.signal
    });

    if (!response.ok) {
      throw new Error(await parseApiError(response));
    }

    const data = await response.json();
    if (!Array.isArray(data?.tags)) {
      throw new Error("La respuesta del servicio no tiene el formato esperado.");
    }

    return data;
  } catch (error) {
    if (error.name === "AbortError") {
      throw new Error("La solicitud excedio el tiempo de espera. Intenta nuevamente.");
    }
    throw error;
  } finally {
    clearTimeout(timeoutId);
  }
}
