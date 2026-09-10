# govalidator

`govalidator` es una librería de validación de datos para Go, inspirada en el estilo
declarativo de frameworks como Laravel: defines una lista de reglas por campo (incluyendo
campos anidados y arreglos, con notación de punto) y la librería valida un `struct` o un
`map[string]any`, regresando los datos ya saneados o una lista de errores legible.

- Reglas built-in: `required`, `email`, `in`, `min`/`max`, comparaciones numéricas y de
  fecha, `unique`/`exists` contra tu propia base de datos, reglas condicionales
  (`required_if`, `required_with`, etc.) y más.
- Soporta objetos anidados y arreglos, tanto de structs como de valores primitivos.
- **Normalizers**: transforma un valor (`lower`, `trim`, `to_int`, etc.) antes de que
  corran las reglas de ese mismo campo.
- Mensajes de error en español, sobreescribibles por campo+regla.

## Instalación

```bash
go get github.com/Nemutagk/govalidator/v2
```

Requiere Go 1.25 o superior (ver `go.mod`).

## Índice

- [Inicio rápido](#inicio-rápido)
- [Conceptos base](#conceptos-base)
- [`ValidateStruct` vs `ValidateRequest`](#validatestruct-vs-validaterequest)
- [Catálogo de reglas](#catálogo-de-reglas)
- [Objetos anidados y arreglos](#objetos-anidados-y-arreglos)
- [Normalizers (pre-proceso de valores)](#normalizers-pre-proceso-de-valores)
- [Mensajes de error personalizados](#mensajes-de-error-personalizados)
- [Manejo de errores](#manejo-de-errores)
- [`unique`, `exists` y `customized`: validaciones contra tu propio código](#unique-exists-y-customized-validaciones-contra-tu-propio-código)
- [Ejemplo completo](#ejemplo-completo)
- [Cosas a tener en cuenta](#cosas-a-tener-en-cuenta)

## Inicio rápido

```go
package main

import (
	"fmt"

	govalidator "github.com/Nemutagk/govalidator/v2"
)

type RegisterRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func main() {
	req := RegisterRequest{Name: "", Email: "no-es-un-correo", Age: 15}

	rules := []govalidator.Input{
		{Name: "name", Rules: []govalidator.Rule{{Name: "required"}}},
		{Name: "email", Rules: []govalidator.Rule{{Name: "required"}, {Name: "email"}}},
		{Name: "age", Rules: []govalidator.Rule{{Name: "required"}, {Name: "greater_than_equal", Options: []string{"18"}}}},
	}

	safe, err := govalidator.ValidateStruct(req, rules, nil, nil)
	if err != nil {
		if ve, ok := govalidator.AsValidationError(err); ok {
			for _, msg := range ve.Errors {
				fmt.Println(msg)
			}
			// Imprime, en orden no garantizado (vienen de un mapa interno):
			//   name: El campo name está vacío
			//   email: El campo no es un correo electrónico válido
			//   age: El campo age debe ser mayor o igual que 18
		}
		return
	}

	fmt.Printf("%+v\n", safe) // req ya validado y saneado, mismo tipo RegisterRequest
}
```

## Conceptos base

Todo se arma alrededor de tres tipos, en el paquete raíz `govalidator`:

```go
type Input struct {
	Name        string       // nombre del campo, admite notación de punto: "payload.mode"
	Rules       []Rule       // reglas de validación para este campo
	Normalizers []Normalizer // transformaciones a aplicar antes que Rules
}

type Rule struct {
	Name    string   // ej. "required", "in", "min"
	Options []string // parámetros de la regla, ej. Options: []string{"18"} para min:18
}

type Normalizer struct {
	Name    string
	Options []string
}
```

Declaras un `[]Input` — uno por cada campo que quieras validar — y se lo pasas a
`ValidateStruct` o `ValidateRequest` junto con el dato a validar.

## `ValidateStruct` vs `ValidateRequest`

- **`ValidateStruct[T any](s T, inputs []Input, customErrors map[string]string, models map[string]func(...)) (T, error)`**
  Recibe un `struct` (usa las etiquetas `json` por default para mapear los nombres de
  campo; puedes usar otra etiqueta con `StructOptions{Tag: "form"}` como último
  argumento opcional). Regresa el mismo tipo `T`, ya saneado.

- **`ValidateRequest(body map[string]any, inputs []Input, customErrors map[string]string, models map[string]func(...)) (map[string]any, error)`**
  Recibe directamente un `map[string]any` (por ejemplo, el resultado de decodificar un
  JSON con `json.Unmarshal` a `map[string]any`, o un `form` ya parseado). Regresa un
  `map[string]any` saneado — solo con las claves que tenían al menos una regla asociada.

Ambas regresan el mismo tipo de error (`govalidator.ValidationError`) cuando la
validación falla; ver [Manejo de errores](#manejo-de-errores).

## Catálogo de reglas

Todas se usan como `Rule{Name: "...", Options: []string{...}}`. `payload` más abajo se
refiere al `map[string]any` completo que se está validando (para `ValidateStruct`, es el
struct convertido a mapa).

| Regla | Options | Qué valida |
|---|---|---|
| `required` | — | El campo existe y no es `nil`/`""`. Si falla, detiene el resto de reglas de ese campo. |
| `sometimes` | — | Si el campo no existe o es `nil`, se saltan las demás reglas de ese campo (no marca error). |
| `nullable` | — | El campo debe existir en el payload (puede ser `nil`). |
| `email` | — | Formato de correo válido (regex). |
| `in` | valores permitidos | El valor debe ser exactamente uno de `Options`. |
| `not_in` | valores prohibidos | El valor no debe estar en `Options`. |
| `equal` | `[valor]` | El valor debe ser igual (`==`) a `Options[0]`. |
| `not_equal` | `[valor]` | El valor debe ser distinto de `Options[0]`. |
| `confirmation` | — | Compara `campo` contra `campo_confirmation` (ej. `password`/`password_confirmation`). |
| `min` | `[n]` | Longitud mínima (string), valor mínimo (int/float) o cantidad mínima de elementos (slice/array). |
| `max` | `[n]` | Igual que `min` pero como tope máximo. |
| `greater_than` / `greater_than_equal` | `[objetivo, layoutPropio?, layoutObjetivo?]` | Compara números o fechas. `objetivo` puede ser un literal o el nombre de otro campo del payload. |
| `less_than` / `less_than_equal` | igual que arriba | Idéntico pero en sentido inverso. |
| `before` / `after` | `[objetivo, formato?]` | Compara fechas contra otro campo, un literal, o las palabras `now`/`today`/`tomorrow`/`yesterday`. |
| `date` | `[layout?]` | Es una fecha válida con el layout dado (default `2006-01-02T15:04:05`). |
| `date_format` | `[layout]` | Igual, pero el layout es obligatorio. |
| `boolean` | — | El valor es literalmente `bool` (no acepta `"true"`/`1`; para eso ver el normalizer `to_bool`). |
| `type` | `[tipo, "nullable"?]` | El tipo de Go del valor (`reflect.TypeOf(...).String()`) coincide con `Options[0]`, ej. `"string"`, `"int"`, `"float64"`. |
| `array` | — | El valor es un slice o array. |
| `ip` | — | Una o varias (separadas por coma) direcciones IPv4/IPv6 válidas. |
| `password` | — | Mínimo 6 caracteres, con al menos un número, una minúscula, una mayúscula y un carácter especial (`$#%&/()!_-`). Acumula todos los errores que falten, no se detiene en el primero. |
| `required_with` | `[otroCampo, valorEsperado?]` | Requerido si `otroCampo` está definido (y, opcionalmente, si vale `valorEsperado`). |
| `required_with_all` | `[campo1, campo2, ...]` | Requerido si **todos** esos campos están definidos. |
| `required_without` | `[otroCampo]` | Requerido si `otroCampo` **no** está definido o está vacío. |
| `required_without_all` | `[campo1, campo2, ...]` | Requerido si **ninguno** de esos campos está definido. |
| `required_if` | `[ruta.con.punto, valorEsperado?]` | Requerido si el nodo en esa ruta (dentro del payload raíz) existe (y, opcionalmente, vale `valorEsperado`). |
| `required_if_all` | `[ruta1, valor1, ruta2, valor2, ...]` | Requerido solo si **todas** esas condiciones se cumplen a la vez (AND). |
| `unique` | `[nombreDeModelo, ...]` | Ver [`unique`, `exists` y `customized`](#unique-exists-y-customized-validaciones-contra-tu-propio-código). |
| `exists` | `[nombreDeModelo, ...]` | Igual, pero valida que SÍ exista (lo opuesto a `unique`). |
| `customized` | `[nombreDeFunción, ...]` | Corre una función de validación 100% tuya. |

## Objetos anidados y arreglos

`Input.Name` soporta notación de punto para llegar a campos anidados o elementos de un
arreglo. Puedes declarar tantos `Input` como necesites sobre el mismo campo padre:

```go
type Payload struct {
	Mode      string `json:"mode"`
	BatchID   string `json:"batch_id"`
}
type Request struct {
	Payload Payload `json:"payload"`
}

rules := []govalidator.Input{
	// regla "bare" (sin punto): valida que el objeto payload exista
	{Name: "payload", Rules: []govalidator.Rule{{Name: "required"}}},
	// reglas con punto: validan sub-campos específicos
	{Name: "payload.mode", Rules: []govalidator.Rule{{Name: "required"}, {Name: "in", Options: []string{"individual", "batch"}}}},
	{Name: "payload.batch_id", Rules: []govalidator.Rule{{Name: "required_if", Options: []string{"payload.mode", "batch"}}}},
}
```

El resultado saneado (`safe.Payload`) solo incluye los sub-campos que tuvieron al menos
una regla — cualquier otro campo del struct/objeto original que no se haya declarado
explícitamente se descarta.

Para arreglos, `*` representa "cada elemento":

```go
type Item struct {
	SKU string `json:"sku"`
	Qty int    `json:"qty"`
}
type Order struct {
	Items []Item `json:"items"`
}

rules := []govalidator.Input{
	{Name: "items", Rules: []govalidator.Rule{{Name: "required"}, {Name: "min", Options: []string{"1"}}}},
	{Name: "items.*.sku", Rules: []govalidator.Rule{{Name: "required"}}},
	{Name: "items.*.qty", Rules: []govalidator.Rule{{Name: "required"}, {Name: "greater_than", Options: []string{"0"}}}},
}
```

También puedes apuntar a un índice específico (`items.0.sku`) en vez de `*`, y arreglos
de valores primitivos (`[]string`, `[]int`, etc.) con `Name: "tags"` + `Rules` aplicadas
directo sobre cada valor.

## Normalizers (pre-proceso de valores)

Un `Normalizer` transforma el valor de un campo **antes** de que corran sus `Rules`.
Todos los `Normalizers` de un campo corren primero, en el orden en que los declaraste;
después corren todas las `Rules` con el valor ya transformado — no se intercalan.

```go
rules := []govalidator.Input{
	{
		Name:        "status",
		Normalizers: []govalidator.Normalizer{{Name: "trim"}, {Name: "lower"}},
		Rules:       []govalidator.Rule{{Name: "in", Options: []string{"active", "inactive"}}},
	},
}
// "  ACTIVE  " pasa la regla "in" y el valor final guardado es "active"
```

Catálogo de normalizers built-in:

| Nombre | Options | Qué hace |
|---|---|---|
| `lower` | — | Convierte a minúsculas. |
| `upper` | — | Convierte a mayúsculas. |
| `trim` | — | Quita espacios en ambos extremos. |
| `ltrim` | — | Quita espacios solo a la izquierda. |
| `rtrim` | — | Quita espacios solo a la derecha. |
| `trim_char` | `[caracteres]` | Quita los caracteres dados (unidos como set) de ambos extremos, ej. `Options: []string{"-"}` convierte `"--id--"` en `"id"`. |
| `capitalize` | — | Pone en mayúscula solo la primera letra (`"juan"` → `"Juan"`). |
| `remove_spaces` | — | Quita **todos** los espacios, no solo los de los extremos. |
| `only_digits` | — | Deja solo los caracteres `0-9` (útil para teléfonos/documentos con formato). |
| `to_int` | — | Convierte string/float numérico a `int`. Si no se puede convertir, deja el valor igual. |
| `to_float` | — | Convierte string/int numérico a `float64`. Igual, no falla si no aplica. |
| `to_bool` | — | Convierte `"true"/"1"/"yes"/"on"` → `true` y `"false"/"0"/"no"/"off"` → `false` (sin distinguir mayúsculas). |

Todos son defensivos: si el tipo del valor no aplica (ej. `lower` sobre un `int`), lo
dejan sin tocar — nunca generan un error de validación por sí mismos. Si necesitas que
un valor mal formado falle la validación, combina el normalizer con una regla downstream
(`type`, `greater_than`, etc.) que si valide el resultado.

Puedes registrar `Normalizers` sobre campos anidados o elementos de arreglo igual que con
`Rules` (`payload.mode`, `items.*.sku`, etc.).

### Efecto secundario importante: mutan el `body` en sitio

Cuando un `Normalizer` transforma un valor, también escribe ese valor de vuelta en el
`body` que se está validando — así, reglas que leen **otro** campo (`confirmation`,
`required_if`, `equal`, etc.) ven la versión ya normalizada, no la cruda. Esto significa
que si llamas `ValidateRequest(miMapa, rules, ...)` directo con tu propio
`map[string]any` y alguna regla usa `Normalizers`, **tu mapa original queda modificado**
después de la llamada. Si necesitas conservar el body original intacto, pásale una copia.

Con `ValidateStruct` esto no aplica: el mapa que se muta es una copia interna generada a
partir de tu struct, nunca tu variable original.

## Mensajes de error personalizados

El tercer argumento de `ValidateStruct`/`ValidateRequest` es un `map[string]string` con
llaves `"campo.regla"`:

```go
customErrors := map[string]string{
	"email.required": "El correo es obligatorio",
	"email.email":    "Ese correo no tiene un formato válido",
	"age.greater_than_equal": "Debes ser mayor de edad",
}

safe, err := govalidator.ValidateStruct(req, rules, customErrors, nil)
```

## Manejo de errores

Cuando falla al menos una regla, ambas funciones regresan un
`govalidator.ValidationError` (implementa `error`). Conviértelo con
`govalidator.AsValidationError`:

```go
_, err := govalidator.ValidateStruct(req, rules, nil, nil)
if err != nil {
	if ve, ok := govalidator.AsValidationError(err); ok {
		ve.Errors  // []string, formato "campo: mensaje"
		ve.Summary // los mismos mensajes unidos con "; "
		ve.Message // "Se encontraron errores en la validación"
	}
}
```

`ValidationErrorResponse` (el tipo que regresa `AsValidationError`) también trae
etiquetas `json` (`message`, `field_errors`, `summary`), listo para responder en un
endpoint HTTP.

## `unique`, `exists` y `customized`: validaciones contra tu propio código

Estas tres reglas no saben nada de bases de datos ni de tu dominio — reciben una función
tuya a través del cuarto argumento (`models`), un
`map[string]func(data any, payload map[string]any, opts *[]string) (bool, string)`:

```go
models := map[string]func(data any, payload map[string]any, opts *[]string) (bool, string){
	"users_by_email": func(data any, payload map[string]any, opts *[]string) (bool, string) {
		email := data.(string)
		existe := miRepositorio.ExisteEmail(email)
		if existe {
			return false, "" // false = la validación NO pasa; "" = usa el mensaje default
		}
		return true, ""
	},
}

rules := []govalidator.Input{
	{Name: "email", Rules: []govalidator.Rule{
		{Name: "unique", Options: []string{"users_by_email"}},
	}},
}

safe, err := govalidator.ValidateStruct(req, rules, nil, models)
```

- `unique:nombre,...` — la función debe regresar `true` cuando el valor **no** existe
  todavía (o sea, cuando es válido registrarlo).
- `exists:nombre,...` — la función debe regresar `true` cuando el valor **sí** existe.
- `customized:nombre,...` — regla de propósito general; regresa `true` cuando el valor
  es válido según tu propia lógica.

En los tres casos, cualquier `Options` después del nombre de la función (`opts[1:]`) se
te pasa como `*[]string` en el tercer parámetro, para que definas parámetros extra sin
tener que hardcodear el nombre de tabla/columna en la librería. Si tu función regresa un
segundo valor (`string`) no vacío, se usa como mensaje de error en vez del genérico.

## Ejemplo completo

```go
type CreateOrderRequest struct {
	CustomerEmail string `json:"customer_email"`
	Status        string `json:"status"`
	Items         []struct {
		SKU string `json:"sku"`
		Qty int    `json:"qty"`
	} `json:"items"`
}

rules := []govalidator.Input{
	{
		Name:        "customer_email",
		Normalizers: []govalidator.Normalizer{{Name: "trim"}, {Name: "lower"}},
		Rules:       []govalidator.Rule{{Name: "required"}, {Name: "email"}},
	},
	{
		Name:        "status",
		Normalizers: []govalidator.Normalizer{{Name: "trim"}, {Name: "lower"}},
		Rules:       []govalidator.Rule{{Name: "required"}, {Name: "in", Options: []string{"pending", "paid"}}},
	},
	{Name: "items", Rules: []govalidator.Rule{{Name: "required"}, {Name: "min", Options: []string{"1"}}}},
	{Name: "items.*.sku", Normalizers: []govalidator.Normalizer{{Name: "upper"}}, Rules: []govalidator.Rule{{Name: "required"}}},
	{Name: "items.*.qty", Rules: []govalidator.Rule{{Name: "required"}, {Name: "greater_than", Options: []string{"0"}}}},
}

safe, err := govalidator.ValidateStruct(req, rules, nil, nil)
if err != nil {
	ve, _ := govalidator.AsValidationError(err)
	// ...responder ve.Errors
}
```

## Cosas a tener en cuenta

- El resultado saneado solo incluye los campos que tuvieron al menos un `Input`
  declarado — cualquier campo del struct/mapa original sin reglas se descarta.
- Cuando `required` (o cualquier regla `required_*`) falla para un campo, el resto de
  reglas de **ese mismo campo** se saltan (no tiene sentido, por ejemplo, correr `email`
  sobre un valor que ya sabemos que no existe).
- `sometimes` sobre un campo padre se propaga automáticamente a sus sub-campos con
  notación de punto: si el padre no existe, no se exige que los hijos existan tampoco.
- Los `Normalizers` mutan el `body` en sitio — ver la nota en la sección de Normalizers
  si validas un `map[string]any` que también uses después de llamar a `ValidateRequest`.
- Puedes usar tipos con nombre para tus campos (ej. `type Status string` con constantes
  como enum) con `ValidateStruct`: el valor se desenvuelve a su primitivo subyacente
  antes de validarse, así que `in`, `equal`, `boolean`, `min`/`max`, etc. funcionan igual
  que con un `string`/`int`/`bool` plano.
- Una regla `campo.*` sobre un arreglo de valores primitivos (`[]string`, `[]int`, o un
  tipo con nombre de esos, no un arreglo de structs/objetos) no valida cada elemento hoy
  — solo funcionan las reglas declaradas directo sobre `campo` (ej. `required`, `min`,
  `array`). Para validar cada elemento de un arreglo necesitas que sus elementos sean
  objetos (`[]structs`), usando `campo.*.subcampo`.
