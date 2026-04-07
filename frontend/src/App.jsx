import { useEffect, useMemo, useState } from "react";
import ImageUploadForm from "./components/ImageUploadForm";
import PreviewCard from "./components/PreviewCard";
import StatusMessage from "./components/StatusMessage";
import TagsList from "./components/TagsList";
import { analyzeImage } from "./services/api";

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

      const result = await analyzeImage(selectedFile);
      setTags(result.tags);
      setInfo(result.tags.length ? "Analisis completado." : "No se detectaron etiquetas.");
    } catch (requestError) {
      setTags([]);
      setInfo("");
      setError(requestError.message || "No fue posible analizar la imagen.");
    } finally {
      setIsLoading(false);
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
