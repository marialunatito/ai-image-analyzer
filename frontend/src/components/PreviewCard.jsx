export default function PreviewCard({ imagePreview, imageName }) {
  if (!imagePreview) {
    return null;
  }

  return (
    <article className="preview-card">
      <h2>Imagen cargada</h2>
      <img src={imagePreview} alt={imageName || "Vista previa de la imagen"} />
      <p>{imageName}</p>
    </article>
  );
}
