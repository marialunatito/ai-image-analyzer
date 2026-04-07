export default function ImageUploadForm({
  isLoading,
  selectedFile,
  onFileChange,
  onSubmit,
  maxSizeMB
}) {
  return (
    <form className="upload-form" onSubmit={onSubmit}>
      <label htmlFor="image-input" className="file-label">
        Seleccionar imagen
      </label>
      <input
        id="image-input"
        name="image"
        type="file"
        accept="image/png,image/jpeg,image/webp,image/jpg"
        onChange={onFileChange}
        disabled={isLoading}
      />
      <p className="helper-text">
        Formatos permitidos: JPG, JPEG, PNG, WEBP. Maximo {maxSizeMB} MB.
      </p>

      <button type="submit" disabled={!selectedFile || isLoading}>
        {isLoading ? "Analizando..." : "Analizar"}
      </button>
    </form>
  );
}
