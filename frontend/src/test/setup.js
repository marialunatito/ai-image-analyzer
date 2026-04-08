import "@testing-library/jest-dom/vitest";

if (!URL.createObjectURL) {
	URL.createObjectURL = () => "blob:test-url";
}

if (!URL.revokeObjectURL) {
	URL.revokeObjectURL = () => {};
}
