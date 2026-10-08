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
- Regla `in_if`: como `in`, pero el catálogo de valores permitidos solo se exige cuando
  otro campo (ruta con notación de punto sobre el payload raíz) es igual a un valor
  dado; si no coincide, no hace nada. Igual que con `required_if`, declarar varios
  `Input` con distintas condiciones sobre el mismo campo da un OR/if-else entre ellas.
- Regla `len`: valida una longitud exacta (`Options[0]`). Para `string` cuenta
  caracteres (no bytes), para slice/array/map cuenta elementos; los demás tipos se
  ignoran. Sin `Options` o con un valor no numérico agrega un error en vez de fallar.
- Regla `uuid`: valida un UUID real de versión 4 a 8 (v7 incluido). `Options`
  opcionales, clasificadas por valor y sin importar el orden: formatos (`canonical` por
  default, `simple`, `braces`, `urn`) y versiones (`4` a `8`). Rechaza nil, max,
  variantes que no son RFC 4122 y UUID repetitivos inventados. No valida fechas.
- `README.md` y `CHANGELOG.md`.

### Changed

- **[Cambio incompatible]** `Input.Parent` deja de ser público. Era contabilidad interna
  (se auto-calculaba al separar un `Name` con notación de punto, para propagar
  `sometimes` a sub-campos) y ningún código dentro de este repo lo asignaba desde
  afuera. Si tu código construía `Input{..., Parent: "..."}` explícitamente, tendrás que
  quitar esa asignación.

- Reglas `min` y `max`: en strings ahora cuentan caracteres en vez de bytes (`"ñ"`
  cuenta como 1, antes como 2), y sin `Options` agregan un error en vez de provocar un
  panic.

### Fixed

- `unique` y `exists` solo aceptaban strings: con un valor numérico o `bool` (por ejemplo
  un id `int64`) respondían `el valor no es válido` sin llamar a tu función. Ahora ese
  valor se convierte a su texto (`int64(123)` llega como `'123'`) y tu función recibe
  siempre un `string`, igual que antes, así que el código existente no cambia. Solo dan
  error `nil`, slices, maps y structs. El mensaje de `unique` para una configuración
  inválida estaba en inglés (`the options is not valid`); ahora es
  `la configuración de conexión no es válida`, igual que en `exists`.
- `before` y `after` sobre un valor numérico solo aceptaban un entero literal como objetivo
  y, si era el nombre de otro campo, lo trataban en silencio como 0. Ahora aceptan, igual
  que `greater_than` y `less_than`, un literal numérico (entero o decimal) o el nombre de
  otro campo numérico, y un objetivo que no sea numérico agrega el error
  `El valor a comparar X no es un número válido` en vez de usarse como 0. El mensaje dice
  `al campo X` cuando el objetivo es un campo.
- `in`, `not_in`, `equal`, `not_equal`, `in_if` y el valor esperado de `required_with`
  comparaban el valor del campo contra `Options` (siempre strings) con `==` entre
  interfaces, así que un valor que no fuera string nunca era igual: `in` y `equal`
  fallaban siempre con un campo numérico o `bool`, y `not_in` y `not_equal` pasaban siempre
  en silencio (`not_in: ['0']` aceptaba el entero 0). Ahora la comparación es texto contra
  texto: strings, `bool`, enteros, `uint` y flotantes (incluidos los tipos con nombre) se
  convierten a su texto antes de compararse, de modo que `5` coincide con `'5'`.
- Los mensajes personalizados de `not_in` y `nullable` solo se aplicaban con las claves
  `campo.notin` y `campo.null`, no con `campo.not_in` y `campo.nullable` (la forma
  `campo.regla` documentada). Ahora funcionan las dos formas; las claves anteriores
  siguen valiendo y la documentada tiene prioridad si se definen ambas.
- `min` y `max` ignoraban en silencio cualquier valor numérico que no fuera `int` ni
  `float64` (`int8`/`int16`/`int32`/`int64`, `uint*`, `float32`): la regla terminaba sin
  error y pasaba. Con un campo `int64` de un struct, `min: 19456` aceptaba `1024`. Ahora
  validan todos los tipos numéricos de Go, incluidos los tipos con nombre. Los enteros se
  comparan como `int64` (sin pasar por `float64`) y un `uint64` mayor que `math.MaxInt64`
  se trata como mayor que cualquier límite.
- `greater_than`, `greater_than_equal`, `less_than` y `less_than_equal` rechazaban con
  `debe ser un número o una fecha válida` los valores `int8`/`int16`/`int32`, `uint*` y
  `float32`, porque solo reconocían `int`, `int64` y `float64`. Un `float32` se convierte
  por su representación decimal corta, así que `float32(0.1)` se compara como `0.1`.
- `before` y `after` solo comparaban como número los valores `int`; un `int64` u otro tipo
  numérico caía en la rama de fechas y fallaba. Ahora usan la misma conversión numérica.
- La regla `ip` provocaba un panic cuando el valor no era un string (por ejemplo un
  número). Ahora ignora los valores que no son string, igual que `nil` y `""`.
- Las reglas `before` y `after` provocaban un panic cuando el valor, o el campo contra el
  que se comparan, no era un string (por ejemplo `bool`, `float64`, un slice o un map).
  Ahora agregan el error de fecha inválida correspondiente.
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
- Un campo de un tipo con nombre distinto de su primitivo subyacente (ej.
  `type StatusType string`, `type Priority int`) siempre fallaba reglas que comparan por
  igualdad o hacen un type assertion sobre un primitivo plano (`in`, `not_in`, `equal`,
  `not_equal`, `boolean`, `min`, `max`, comparaciones numéricas), sin importar el valor,
  porque `ValidateStruct` conservaba el tipo con nombre al convertir el struct a
  `map[string]any` — y en Go, dos interfaces solo son iguales si su tipo dinámico también
  coincide. Ahora esos valores se "desenvuelven" a su primitivo subyacente
  (`string`, `bool`, `int`, `int8`...`uint64`, `float32`/`64`) antes de validarse.

### Tests

- Pruebas que recorren todos los tipos numéricos de Go en `min`, `max`, `greater_than*`,
  `less_than*`, `before` y `after` (`validate/numeric_types_test.go`), más el caso de un
  campo `int64` validado a través de `ValidateStruct`.
- Pruebas de regresión para los panics de `ip`, `before` y `after`, y para `len`, `uuid`
  y los cambios de `min` y `max`.
- Suite de pruebas de regresión a nivel raíz (`validate_test.go`, `struct_test.go`) para
  los tres fixes anteriores y para el nuevo sistema de `Normalizers`.
