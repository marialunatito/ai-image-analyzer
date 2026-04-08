import { afterEach, describe, expect, it, vi } from "vitest";
import { ApiError, analyzeImage } from "./api";

describe("analyzeImage", () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("envia la imagen y retorna tags", async () => {
    const fetchMock = vi.spyOn(globalThis, "fetch").mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => ({ tags: [{ label: "Dog", confidence: 0.98 }] })
    });

    const file = new File(["img"], "dog.png", { type: "image/png" });
    const result = await analyzeImage(file);

    expect(result.tags).toHaveLength(1);
    expect(fetchMock).toHaveBeenCalledTimes(1);
  });

  it("lanza ApiError para errores HTTP", async () => {
    vi.spyOn(globalThis, "fetch").mockResolvedValue({
      ok: false,
      status: 415,
      json: async () => ({ code: "INVALID_IMAGE", message: "formato invalido" })
    });

    const file = new File(["img"], "bad.txt", { type: "text/plain" });

    await expect(analyzeImage(file)).rejects.toMatchObject({
      name: "ApiError",
      code: "INVALID_IMAGE",
      message: "formato invalido"
    });
  });

  it("retorna error cancelado cuando signal externo aborta", async () => {
    vi.spyOn(globalThis, "fetch").mockImplementation(
      (_, { signal }) =>
        new Promise((resolve, reject) => {
          signal.addEventListener("abort", () => reject(new DOMException("Aborted", "AbortError")));
          setTimeout(
            () =>
              resolve({
                ok: true,
                status: 200,
                json: async () => ({ tags: [] })
              }),
            50
          );
        })
    );

    const controller = new AbortController();
    const file = new File(["img"], "dog.png", { type: "image/png" });
    const request = analyzeImage(file, { signal: controller.signal });

    controller.abort();

    await expect(request).rejects.toBeInstanceOf(ApiError);
    await expect(request).rejects.toMatchObject({ code: "CANCELED" });
  });
});
