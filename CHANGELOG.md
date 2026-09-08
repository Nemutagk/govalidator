# Changelog

Todos los cambios notables de este proyecto se documentan en este archivo.

El formato sigue [Keep a Changelog](https://keepachangelog.com/es-ES/1.1.0/) y este
proyecto intenta seguir [Semantic Versioning](https://semver.org/lang/es/).

> Este archivo se creó a partir de la versión `v2.15.1`. El historial anterior (desde
> `v1.0.0`) no se documentó retroactivamente aquí; consulta `git log` y `git tag` para
> revisarlo.

## [Unreleased]

### Added

- `Normalizer` / `Input.Normalizers`: nuevo campo, independiente de `Rules`, para
  transformar el valor de un campo (`lower`, `upper`, `trim`, `ltrim`, `rtrim`,
  `trim_char`, `capitalize`, `remove_spaces`, `only_digits`, `to_int`, `to_float`,
  `to_bool`) antes de que corran sus reglas de validación. Todos los normalizers de un
  campo corren primero, en el orden declarado, y el resultado se escribe de vuelta en el
  body de validación para que reglas que leen otro campo (`confirmation`,
  `required_if`, etc.) vean el valor ya normalizado.
- Nuevo paquete `normalize/` con la implementación de cada normalizer built-in y sus
  pruebas unitarias.
- `README.md` y `CHANGELOG.md`.

### Changed

- **[Cambio incompatible]** `Input.Parent` deja de ser público. Era contabilidad interna
  (se auto-calculaba al separar un `Name` con notación de punto, para propagar
  `sometimes` a sub-campos) y ningún código dentro de este repo lo asignaba desde
  afuera. Si tu código construía `Input{..., Parent: "..."}` explícitamente, tendrás que
  quitar esa asignación.

### Fixed

- Una regla sin punto (ej. `payload` o `items`) sobre un campo cuyo valor es un
  struct/objeto anidado (`map[string]any`) siempre fallaba, sin importar qué tan válido
  fuera el contenido, con un mensaje de error duplicado
  (`payload.payload: El campo payload no está definido`). La regla se aplicaba contra sí
  misma en vez de contra el objeto padre.
- Cuando esa misma regla sin punto coexistía con reglas con punto sobre sub-campos del
  mismo objeto (ej. `payload` junto con `payload.mode`), el resultado validado terminaba
  sobreescrito con el valor crudo sin filtrar, reintroduciendo campos que ninguna regla
  había cubierto. Esto ya pasaba también con arreglos (`items` + `items.*.mode`); ahora
  el comportamiento es consistente para ambos casos.

### Tests

- Suite de pruebas de regresión a nivel raíz (`validate_test.go`) para los dos fixes
  anteriores y para el nuevo sistema de `Normalizers`.
