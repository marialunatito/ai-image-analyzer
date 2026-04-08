const DEFAULT_TIMEOUT_MS = 20000;

export class ApiError extends Error {
  constructor(code, message, status) {
    super(message);
    this.name = "ApiError";
    this.code = code;
    this.status = status;
  }
}

function buildApiUrl(path) {
  const base = import.meta.env.VITE_API_BASE_URL || "";
  return `${base}${path}`;
}

async function parseApiError(response) {
  const fallbackMessage = `Error del servidor (${response.status}).`;
  try {
    const data = await response.json();
    return {
      message: data?.message || fallbackMessage,
      code: data?.code || "HTTP_ERROR"
    };
  } catch {
    return { message: fallbackMessage, code: "HTTP_ERROR" };
  }
}

export async function analyzeImage(file, options = {}) {
  const { signal: externalSignal } = options;
  const controller = new AbortController();

  if (externalSignal) {
    externalSignal.addEventListener("abort", () => controller.abort(), { once: true });
  }

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
      const parsedError = await parseApiError(response);
      throw new ApiError(parsedError.code, parsedError.message, response.status);
    }

    const data = await response.json();
    if (!Array.isArray(data?.tags)) {
      throw new ApiError(
        "INVALID_RESPONSE",
        "La respuesta del servicio no tiene el formato esperado.",
        response.status
      );
    }

    return data;
  } catch (error) {
    if (error.name === "AbortError") {
      if (externalSignal?.aborted) {
        throw new ApiError("CANCELED", "Analisis cancelado por el usuario.");
      }

      throw new ApiError("TIMEOUT", "La solicitud excedio el tiempo de espera. Intenta nuevamente.");
    }

    if (error instanceof ApiError) {
      throw error;
    }

    throw new ApiError("NETWORK_ERROR", "No se pudo conectar con el backend.");
  } finally {
    clearTimeout(timeoutId);
  }
}
