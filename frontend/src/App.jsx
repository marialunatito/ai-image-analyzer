import { useEffect, useMemo, useRef, useState } from "react";
import ImageUploadForm from "./components/ImageUploadForm";
import PreviewCard from "./components/PreviewCard";
import StatusMessage from "./components/StatusMessage";
import TagsList from "./components/TagsList";
import { ApiError, analyzeImage } from "./services/api";

const DEFAULT_MAX_SIZE_MB = Number(import.meta.env.VITE_MAX_IMAGE_SIZE_MB || 5);
const ALLOWED_MIME_TYPES = [
  "image/png",
  "image/jpeg",
  "image/jpg",
  "image/webp"
];

function validateImage(file, maxSizeMB) {
  if (!file) {
    return "Selecciona una imagen para continuar.";
  }

  if (!ALLOWED_MIME_TYPES.includes(file.type)) {
    return "Formato no permitido. Usa JPG, JPEG, PNG o WEBP.";
  }

  const maxSizeBytes = maxSizeMB * 1024 * 1024;
  if (file.size > maxSizeBytes) {
    return `La imagen excede el limite de ${maxSizeMB} MB.`;
  }

  return null;
}

export default function App() {
  const abortControllerRef = useRef(null);
  const [selectedFile, setSelectedFile] = useState(null);
  const [imagePreview, setImagePreview] = useState("");
  const [tags, setTags] = useState([]);
  const [error, setError] = useState("");
  const [info, setInfo] = useState("Carga una imagen y analizala con IA.");
  const [isLoading, setIsLoading] = useState(false);

  const maxSizeMB = useMemo(
    () => (Number.isFinite(DEFAULT_MAX_SIZE_MB) ? DEFAULT_MAX_SIZE_MB : 5),
    []
  );

  useEffect(() => {
    if (!selectedFile) {
      setImagePreview("");
      return;
    }

    const objectUrl = URL.createObjectURL(selectedFile);
    setImagePreview(objectUrl);

    return () => {
      URL.revokeObjectURL(objectUrl);
    };
  }, [selectedFile]);

  function handleFileChange(event) {
    const file = event.target.files?.[0];
    setTags([]);
    setError("");
    setInfo("Imagen lista para analizar.");

    const validationError = validateImage(file, maxSizeMB);
    if (validationError) {
      setSelectedFile(null);
      setInfo("");
      setError(validationError);
      return;
    }

    setSelectedFile(file);
  }

  function handleCancel() {
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
    }
  }

  async function handleSubmit(event) {
    event.preventDefault();

    const validationError = validateImage(selectedFile, maxSizeMB);
    if (validationError) {
      setError(validationError);
      setInfo("");
      return;
    }

    try {
      setIsLoading(true);
      setError("");
      setInfo("Analizando imagen, por favor espera...");
      abortControllerRef.current = new AbortController();

      const result = await analyzeImage(selectedFile, {
        signal: abortControllerRef.current.signal
      });
      setTags(result.tags);
      setInfo(result.tags.length ? "Analisis completado." : "No se detectaron etiquetas.");
    } catch (requestError) {
      setTags([]);
      setInfo("");
      setError(resolveErrorMessage(requestError));
    } finally {
      setIsLoading(false);
      abortControllerRef.current = null;
    }
  }

  return (
    <main className="page">
      <header className="hero">
        <p className="kicker">Prueba tecnica Full-Stack</p>
        <h1>Analizador Inteligente de Imagenes</h1>
        <p>
          Sube una imagen y revisa etiquetas generadas por IA con su nivel de confianza.
        </p>
      </header>

      <section className="panel">
        <ImageUploadForm
          isLoading={isLoading}
          selectedFile={selectedFile}
          onFileChange={handleFileChange}
          onSubmit={handleSubmit}
          onCancel={handleCancel}
          maxSizeMB={maxSizeMB}
        />

        {isLoading && <div className="spinner" aria-label="Cargando" />}

        <StatusMessage variant="error">{error}</StatusMessage>
        <StatusMessage variant="info">{info}</StatusMessage>

        <div className="results-grid">
          <PreviewCard imagePreview={imagePreview} imageName={selectedFile?.name} />
          <TagsList tags={tags} />
        </div>
      </section>
    </main>
  );
}

function resolveErrorMessage(error) {
  if (!(error instanceof ApiError)) {
    return error?.message || "No fue posible analizar la imagen.";
  }

  switch (error.code) {
    case "CANCELED":
      return "Analisis cancelado.";
    case "TIMEOUT":
      return "La IA tardo demasiado. Intenta con otra imagen o vuelve a intentar.";
    case "INVALID_IMAGE":
    case "INVALID_REQUEST":
    case "PAYLOAD_TOO_LARGE":
      return error.message;
    case "PROVIDER_ERROR":
      return "El proveedor de IA no esta disponible en este momento.";
    default:
      return error.message || "No fue posible analizar la imagen.";
  }
}
