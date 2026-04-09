# backend/test

Esta carpeta se usa para pruebas E2E e integracion.

Estado actual:

- E2E: test/e2e/analyze_e2e_test.go
- Placeholder historico: test/main_test.go

Convencion del proyecto:

- Los tests unitarios viven junto al codigo del paquet
e.
- Ejemplos: internal/handler/analyze_test.go, internal/usecase/analyze_test.go, internal/service/ia_test.go.

Comandos utiles:

- Ejecutar solo E2E: go test ./test/e2e -count=1
- Ejecutar todo backend: go test ./... -count=1
