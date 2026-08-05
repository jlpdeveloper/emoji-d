# ADR 002: Phase 1 Uses a Curated 32-Emoji Alphabet

## Status

Accepted

## Context

`emojiid.ID` needs a deterministic validation rule for Phase 1.

ADR 001 established that Phase 1 validates IDs as non-empty UTF-8 strings made only from supported single-rune emoji. The remaining question is how large the supported emoji alphabet should be.

Unicode contains many emoji, but accepting every single-codepoint emoji would make the initial package harder to reason about. It would also include symbols that may be visually ambiguous, presentation-sensitive, rude in unexpected ways, hard to type, or simply not useful as identifiers.

The project should stay small, readable, and manually reviewable.

## Decision

Phase 1 will use a private, curated alphabet of exactly 32 supported emoji runes.

An `emojiid.ID` is valid when every rune in the ID exists in this 32-symbol alphabet.

This gives the package a simple mental model:

```text
emojiid IDs are base-32 strings, where each digit is an emoji
```

Examples:

```text
🐱
🦆
🐱🦆
🦆🦆
🌮🔥
```

The alphabet is curated manually. It should favor emoji that are:

- single-rune
- visually distinct
- common on major platforms
- reasonably easy to type or recognize
- fun enough for the project

The alphabet should avoid:

- flags
- skin tone modifiers
- zero-width-joiner sequences
- variation-selector-dependent symbols
- near-duplicate emoji
- text-like symbols that may render inconsistently

## Consequences

### Positive

- The supported set is easy to audit.
- The implementation is small.
- Validation is deterministic.
- The test suite can assert the alphabet size is exactly 32.
- The package has a memorable base-32 identity model.
- Multi-emoji IDs naturally expand the available ID space.

Approximate capacity:

```text
1 emoji  = 32 IDs
2 emoji  = 1,024 IDs
3 emoji  = 32,768 IDs
4 emoji  = 1,048,576 IDs
5 emoji  = 33,554,432 IDs
6 emoji  = 1,073,741,824 IDs
```

### Negative

- Many valid Unicode emoji are not supported.
- Consumers cannot customize the alphabet in Phase 1.
- The default alphabet is subjective.
- Some users may dislike specific emoji included in the curated set.

## Alternatives Considered

### Accept All Single-Codepoint Emoji

Rejected for Phase 1.

This would make the package feel more complete, but it also pulls in a large and uneven support surface. The project benefits more from a small, explicit alphabet than from broad Unicode coverage.

### Use 64 Emoji

Deferred.

A 64-symbol alphabet would provide a neat base-64 model and a much larger ID space per character. However, 32 is easier to curate manually and is already large enough for a Phase 1 caller-supplied ID system.

### Let Consumers Provide Their Own Alphabet

Deferred to Phase 2.

Custom alphabets are useful, but they introduce API design questions around validators, defaults, generated code, and package-level behavior.

## Future Work

Potential Phase 2 additions:

- exported alphabets
- custom validators
- 64-symbol default alphabet
- generated candidate lists from Unicode data
- tooling to verify that each alphabet entry is single-rune
- documentation showing supported emoji

## Summary

Phase 1 chooses a private 32-emoji alphabet because it keeps validation deterministic, small, and fun.

The package does not try to support all emoji. It supports a curated base-32 emoji alphabet and treats IDs as non-empty strings made from those symbols.
