# emojiid Design Doc

## Purpose

`emojiid` is a small Go package for using emoji strings as intentionally silly, but technically valid, PostgreSQL primary keys.

The project is meant to be fun first. It should stay small enough to build by hand while still teaching real engineering concepts:

- custom Go types for domain values
- `database/sql` integration
- `sqlc` type overrides
- PostgreSQL domains over `text`
- JSON and text marshaling
- Unicode string handling in Go
- delegating uniqueness to database constraints

The main idea is simple: callers bring their own emoji ID, and PostgreSQL rejects duplicates through the primary key constraint.

## What We Are Building

Phase 1 builds a text-backed package named `emojiid`.

The package exposes one concrete type:

```text
emojiid.ID
```

`ID` wraps a private string field. It represents a non-null emoji-flavored identifier that can be passed into normal database code.

Conceptually:

```text
type ID struct {
  value string
}
```

The actual implementation should keep the field private so the package controls validation.

## Storage Model

PostgreSQL stores the value as `text`.

Recommended database shape:

```text
create domain emoji_id as text;
```

Then tables can use:

```text
id emoji_id primary key
```

The domain gives the schema a meaningful type name while keeping the underlying storage simple. Indexes, joins, foreign keys, and primary-key uniqueness all work like normal text.

## Usage Model

The consumer supplies IDs explicitly.

Example flow:

```text
id := emojiid.Must("🐱")
CreateObject(ctx, id, "hello world")
```

The package does not generate IDs in Phase 1. It only validates, stores, scans, and serializes them.

If the caller reuses an existing ID, PostgreSQL rejects the insert with a primary-key violation. That is expected behavior and should not be hidden by the package.

## Phase 1 Guarantees

`emojiid.ID` should guarantee:

- the value is non-empty
- the value is valid UTF-8
- the value passes a deliberately small emoji-ish validation rule
- it can be stored in a PostgreSQL `text` or domain-over-`text` column
- it can be scanned from database query results
- it can be used as a database query argument
- it can be marshaled and unmarshaled as text
- it can be marshaled and unmarshaled as JSON
- it has a useful string representation

## Phase 1 Non-Goals

Phase 1 should not attempt:

- nullable IDs
- automatic ID generation
- bytea-backed IDs
- custom SQL builders
- custom ordering
- complete Unicode emoji validation
- database-side emoji validation
- semantic emoji ranges
- collision prevention outside the database
- multi-database support

This boundary keeps the project light and still useful.

## Validation

Validation should be intentionally modest in Phase 1.

A reasonable first rule:

- reject empty strings
- reject invalid UTF-8
- reject obvious plain ASCII identifiers
- accept values that are plausibly emoji

The package does not need to prove that a string is exactly one Unicode emoji grapheme cluster.

That distinction can become complicated because many visible emoji are composed from multiple Unicode code points. For this project, the better model is:

```text
An ID is a caller-provided emoji-flavored text identifier.
```

This allows values like:

```text
🐱
🦆
🦆🦆
🐱🦆
```

The consumer remains responsible for choosing IDs that make sense for their application.

## Required Interfaces

### database/sql.Scanner

Purpose: lets database query results scan into `emojiid.ID`.

Required method:

```text
Scan(value any) error
```

Expected behavior:

- accept `string`
- accept `[]byte` if returned by the driver
- reject `nil`, because Phase 1 is non-nullable
- reject unsupported input types
- validate the scanned value before storing it

### database/sql/driver.Valuer

Purpose: lets `emojiid.ID` be passed as a database query argument.

Required method:

```text
Value() (driver.Value, error)
```

Expected behavior:

- return the underlying string
- return an error if the ID is invalid
- never return `nil` for a valid `ID`

### encoding.TextMarshaler

Purpose: lets the ID encode as plain text.

Required method:

```text
MarshalText() ([]byte, error)
```

Expected behavior:

- return the underlying string as UTF-8 bytes
- return an error if the ID is invalid

### encoding.TextUnmarshaler

Purpose: lets the ID decode from plain text.

Required method:

```text
UnmarshalText(text []byte) error
```

Expected behavior:

- parse the bytes as UTF-8 text
- validate the resulting string
- store it only if valid

### encoding/json.Marshaler

Purpose: lets the ID appear as a JSON string.

Required method:

```text
MarshalJSON() ([]byte, error)
```

Expected behavior:

- encode `emojiid.ID` as a JSON string, such as `"🐱"`
- return an error if the ID is invalid

### encoding/json.Unmarshaler

Purpose: lets the ID parse from a JSON string.

Required method:

```text
UnmarshalJSON(data []byte) error
```

Expected behavior:

- accept JSON strings
- reject `null` in Phase 1
- reject objects, arrays, numbers, and booleans
- validate the decoded string before storing it

### fmt.Stringer

Purpose: makes the ID readable in logs, errors, and debugging.

Required method:

```text
String() string
```

Expected behavior:

- return the underlying emoji string
- for an invalid zero value, either return an empty string or a clear debug sentinel

Returning an empty string is idiomatic but less obvious in debugging. Returning a sentinel is more visible but can surprise consumers. This should be decided deliberately.

## Public API

Minimal Phase 1 API:

```text
New(value string) (ID, error)
Must(value string) ID
Validate(value string) error
```

Potential utility methods:

```text
String() string
IsValid() bool
```

`New` is for normal application paths.

`Must` is for tests, seed data, examples, and places where a literal should be known-valid.

`Validate` is useful for callers that want to check input before constructing an `ID`.

## sqlc Integration

The intended sqlc setup is:

1. Define a PostgreSQL domain over `text`.
2. Use that domain for ID columns.
3. Configure sqlc to map that database type to `emojiid.ID`.

Conceptual override:

```text
db_type: emoji_id
go_type: github.com/yourname/emojiid.ID
```

After that, generated sqlc structs and methods can use `emojiid.ID` directly.

The SQL remains normal:

```text
insert into objects (id, body) values ($1, $2)
select * from objects where id = $1
```

The Go type handles binding and scanning through the database interfaces.

## PostgreSQL Behavior

Uniqueness is handled by the primary key.

Ordering is handled by PostgreSQL text comparison and the column/database collation. The package should not promise human-meaningful emoji ordering.

Queries like this can work:

```text
where id between $1 and $2
```

But the meaning is:

```text
between these text values according to PostgreSQL ordering
```

not:

```text
between these emoji semantically
```

## Tests To Write

Core constructor tests:

- valid single emoji
- valid repeated emoji
- valid mixed emoji
- empty string rejected
- plain ASCII rejected if that is part of validation
- invalid UTF-8 rejected

Database interface tests:

- `Value` returns a string
- `Value` rejects an invalid zero value
- `Scan` accepts `string`
- `Scan` accepts `[]byte`
- `Scan` rejects `nil`
- `Scan` rejects unsupported types
- `Scan` rejects invalid text

Encoding tests:

- text marshal round trip
- text unmarshal validates input
- JSON marshal emits a string
- JSON unmarshal accepts a string
- JSON unmarshal rejects `null`
- JSON unmarshal rejects non-string JSON

General behavior tests:

- `String` returns expected display value
- zero value behavior is documented and tested
- `Must` panics on invalid input

## Phase 2 Ideas

Possible future additions:

- `NullID`
- fixed emoji alphabets
- ID generation
- bytea-backed `ByteID`
- stricter Unicode emoji validation
- grapheme-cluster-aware validation
- custom comparison helpers
- Postgres check constraints
- sqlc example project
- pgx-specific examples
- README badges, CI, and fuzz tests

These should stay out of Phase 1 unless they become necessary.

## Summary

Phase 1 is intentionally small:

```text
PostgreSQL text domain + sqlc override + Go ID wrapper
```

The database owns uniqueness. The caller owns choosing IDs. The package owns validation and interface compatibility.

That keeps the project silly, teachable, and finishable.
