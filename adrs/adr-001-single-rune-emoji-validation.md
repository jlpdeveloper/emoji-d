# ADR 001: Phase 1 Emoji Validation Uses a Private Single-Rune Alphabet

## Status

Accepted

## Context

`emojiid` is a small Go package for using emoji-flavored text values as PostgreSQL primary keys through normal `database/sql` and `sqlc` workflows.

The package needs a validation rule for `ID` values. A fully correct Unicode emoji validator is more complex than the project needs for Phase 1 because many visible emoji are not a single Unicode code point.

Examples of emoji sequences that complicate validation:

- flags, such as regional-indicator pairs
- skin tone modifiers
- gender and profession variants
- family emoji
- emoji joined with zero-width joiners
- variation-selector forms

These sequences can also create practical lookup and display confusion. Users may choose visually similar emoji with different underlying code point sequences depending on keyboard, platform, or preference.

The project is intentionally silly, but the API should still have a clear and deterministic contract.

## Decision

For Phase 1, `emojiid.ID` validation will use a private allowed set of single-rune emoji.

An `ID` is valid when:

- the value is non-empty
- the value is valid UTF-8
- every rune in the value exists in the package's private allowed emoji set

This means an ID is modeled as:

```text
a non-empty sequence of supported emoji runes
```

Examples:

```text
🐱      valid
🦆      valid
🐱🦆    valid
🦆🦆    valid
hello   invalid
```

Emoji sequences are intentionally out of scope for Phase 1.

Examples that may be rejected:

```text
👍🏽
🇺🇸
👨‍💻
🏳️‍🌈
```

## Consequences

### Positive

- Validation is deterministic.
- The implementation stays small.
- Tests are straightforward.
- Plain ASCII identifiers are naturally rejected.
- Repeated emoji IDs work without special handling.
- Mixed emoji IDs work without special handling.
- The package avoids claiming full Unicode emoji correctness.
- The supported ID space is controlled by the package.

### Negative

- Some visually valid emoji will be rejected.
- Flags, skin tones, and joined emoji are not supported in Phase 1.
- The package's supported emoji set may feel arbitrary until documented.
- Consumers cannot customize the allowed emoji alphabet yet.

## Alternatives Considered

### Full Unicode Emoji Validation

This would attempt to validate all current Unicode emoji, including grapheme clusters and zero-width-joiner sequences.

Rejected for Phase 1 because it adds complexity that is not necessary for the learning goals or the initial package contract.

### Broad Unicode Range Checks

This would validate runes by checking whether they fall inside rough Unicode emoji ranges.

Rejected because emoji are spread across Unicode, and range checks can both reject valid emoji and accept unrelated symbols.

### Any Non-Empty UTF-8 String

This would treat the package as a thin text wrapper with almost no emoji-specific validation.

Rejected because the package should enforce at least one meaningful project invariant: IDs should be made only from supported emoji symbols.

### Public Custom Alphabet

This would let consumers provide their own allowed emoji set.

Deferred to Phase 2. It is useful, but it adds API design questions around validators, package defaults, and sqlc-generated code usage.

## Future Work

Possible Phase 2 additions:

- exported alphabets
- custom validators
- configurable allowed emoji sets
- support for selected emoji sequences
- nullable `NullID`
- ID generation from an alphabet
- stricter docs around ordering and PostgreSQL collation behavior

## Summary

Phase 1 chooses a small, explicit contract:

```text
emojiid.ID = non-empty valid UTF-8 made only from supported single-rune emoji
```

This keeps the project fun, teachable, and finishable while still creating a real domain type with useful database and encoding behavior.
