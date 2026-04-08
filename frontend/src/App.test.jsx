import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import App from "./App";
import * as api from "./services/api";

describe("App", () => {
  it("muestra error cuando archivo no es imagen", () => {
    render(<App />);

    const input = screen.getByLabelText(/seleccionar imagen/i);
    const file = new File(["text"], "file.txt", { type: "text/plain" });

    fireEvent.change(input, { target: { files: [file] } });

    expect(screen.getByText(/formato no permitido/i)).toBeInTheDocument();
  });

  it("analiza una imagen valida y renderiza tags", async () => {
    vi.spyOn(api, "analyzeImage").mockResolvedValue({
      tags: [{ label: "Perro", confidence: 0.9 }]
    });

    render(<App />);

    const input = screen.getByLabelText(/seleccionar imagen/i);
    const file = new File(["img"], "dog.png", { type: "image/png" });
    fireEvent.change(input, { target: { files: [file] } });

    fireEvent.click(screen.getByRole("button", { name: /analizar/i }));

    await waitFor(() => {
      expect(screen.getByText("Perro")).toBeInTheDocument();
    });
  });

  it("permite cancelar analisis", async () => {
    let resolver;
    vi.spyOn(api, "analyzeImage").mockImplementation(
      () =>
        new Promise((_, reject) => {
          resolver = () => reject(new api.ApiError("CANCELED", "Analisis cancelado."));
        })
    );

    render(<App />);

    const input = screen.getByLabelText(/seleccionar imagen/i);
    const file = new File(["img"], "dog.png", { type: "image/png" });
    fireEvent.change(input, { target: { files: [file] } });

    fireEvent.click(screen.getByRole("button", { name: /analizar/i }));
    expect(screen.getByRole("button", { name: /cancelar/i })).toBeInTheDocument();

    resolver();

    await waitFor(() => {
      expect(screen.getByText(/analisis cancelado/i)).toBeInTheDocument();
    });
  });
});
