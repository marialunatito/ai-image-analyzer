import { toPercent } from "../utils/format";

export default function TagsList({ tags }) {
  if (!tags?.length) {
    return null;
  }

  return (
    <section className="tags-card">
      <h2>Etiquetas detectadas</h2>
      <ul>
        {tags.map((tag, index) => (
          <li key={`${tag.label}-${index}`}>
            <div>
              <strong>{tag.label}</strong>
              <span>{toPercent(tag.confidence)}</span>
            </div>
            <progress max="1" value={Number(tag.confidence) || 0} />
          </li>
        ))}
      </ul>
    </section>
  );
}
