# Emoji-d

Have you ever looked at an `int`, squinted at a UUID, and thought: this primary key has absolutely no theatrical value?

Good news. The database did not ask for dignity.

`emojid` is a tiny Go package for using curated emoji strings as PostgreSQL IDs. It wraps a private string, validates that the value is made only from a small approved emoji alphabet, and implements the boring interfaces required to work with `database/sql`, JSON, text encoding, and `sqlc`.

This is not a UUID replacement. This is a small act of database whimsy with just enough engineering discipline to avoid becoming a production incident with merch.

## What It Does

`emojid.ID` is a non-null, text-backed ID type for people who believe `🐱` deserves the same referential integrity as `42`.

It is valid when:

- it is not empty
- it is valid UTF-8
- every rune is in the curated emoji alphabet

That means these are valid:

```text
🐱
🦆
🦆🦆
🌮🔥
🍕🐢👀
```

And these are not:

```text
hello
123
🐱abc
🇺🇸
👍🏽
```

The package intentionally does not support flags, skin-tone variants, zero-width-joiner emoji, or whatever Unicode machinery is required to make a family emoji hold hands consistently across three operating systems and a terminal emulator.

## Why

Because sometimes you need to remember software can still be fun, and sometimes that means telling PostgreSQL that `🦆` is a serious business identifier.

Also, this is a surprisingly decent little learning project:

- custom Go domain types
- `database/sql.Scanner`
- `driver.Valuer`
- JSON marshaling
- text marshaling
- UTF-8 validation
- PostgreSQL domains
- `sqlc` type overrides
- letting the database enforce uniqueness like a responsible adult

## AI Disclosure

This is not a vibe-coded project.

The core implementation was written deliberately, by a human, with tests and docs getting some AI assistance. AI also served as the project's sentient rubber duck: patient, occasionally useful, and strangely willing to discuss whether `🌮` is a better primary key digit than `⭐`.

## Database Shape

Recommended PostgreSQL setup:

```sql
create domain emoji_id as text;

create table objects (
  id emoji_id primary key,
  body text not null
);
```

The domain is still backed by `text`, so PostgreSQL can index it, compare it, join on it, and reject duplicate primary keys without needing to understand the bit.

## sqlc Usage

Use normal SQL:

```sql
-- name: CreateObject :one
insert into objects (id, body)
values ($1, $2)
returning id, body;
```

Then configure `sqlc` to map the database type to your Go type:

```yaml
overrides:
  - db_type: "emoji_id"
    go_type:
      import: "github.com/you/emojid"
      type: "ID"
```

Now generated code can accept and return `emojid.ID` directly.

The caller brings the ID:

```text
CreateObject(ctx, emojid.Must("🐱"), "hello world")
```

If `🐱` is already taken, PostgreSQL rejects the insert. As it should. There can only be one cat. The database has spoken.

## Interfaces

`ID` is designed to satisfy:

```text
database/sql.Scanner
database/sql/driver.Valuer
encoding.TextMarshaler
encoding.TextUnmarshaler
encoding/json.Marshaler
encoding/json.Unmarshaler
fmt.Stringer
```

The important behavior:

- database values scan into `ID`
- query args bind as text
- JSON encodes as `"🐱"`, not `{"ID":"🐱"}`
- text encoding uses the raw emoji string
- invalid zero values should not quietly sneak into the database

## Curated Alphabet

Phase 1 uses exactly 32 supported emoji. Think of it as base-32, except the digits have attitude and one of them is cheese.

```go
package emojid

// validRunes is a curated list of emojis that are allowed for use in IDs.
// Source: https://unicode.org/emoji/charts/full-emoji-list.html
var validRunes = map[rune]struct{}{
	'😀': {},
	'🚀': {},
	'🔥': {},
	'🍕': {},
	'😅': {},
	'😐': {},
	'🙄': {},
	'🤮': {},
	'🤠': {},
	'🌮': {},
	//10
	'😎': {},
	'💩': {},
	'👻': {},
	'😺': {},
	'🍝': {},
	'💯': {},
	'🖕': {},
	'👍': {},
	'💪': {},
	'👀': {},
	//20
	'🦆': {},
	'🦉': {},
	'🐔': {},
	'🐸': {},
	'🐢': {},
	'🐌': {},
	'🌭': {},
	'🌲': {},
	'🍄': {},
	'🍌': {},
	//30
	'🍆': {},
	'🧀': {},
}
```

Capacity, if you want to pretend this is a numbering system:

```text
1 emoji  = 32 possible IDs
2 emoji  = 1,024 possible IDs
3 emoji  = 32,768 possible IDs
4 emoji  = 1,048,576 possible IDs
5 emoji  = 33,554,432 possible IDs
6 emoji  = 1,073,741,824 possible IDs
```

But Phase 1 does not generate IDs. You pick them. The database tells you when your vibes collide. This is called architecture.

## Non-Goals

This package does not try to:

- generate IDs
- guarantee uniqueness before insert
- support nullable IDs
- support every Unicode emoji
- support skin tones or flags
- define semantic emoji ordering
- replace migrations, query builders, or your sense of judgment

Range queries may work because PostgreSQL can compare text and nobody can stop you:

```sql
where id between $1 and $2
```

But that means PostgreSQL text ordering, not “all emoji between duck and taco in the grand taxonomy of nonsense.” Please do not cite this package in court.

## Philosophy

The database owns uniqueness.

The caller owns choosing the ID.

`emojid` owns one tiny invariant:

```text
This string is made only of approved emoji runes.
```

That is enough. Anything more would risk turning a joke into enterprise software, and nobody needs an Emoji Identity Governance Committee.
