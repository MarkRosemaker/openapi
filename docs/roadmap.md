# Roadmap

Work not yet done. An entry is deleted once it is.

## Normal and strict validation

`Validate` rejects what the specifications forbid, a MUST. Real documents
also break what they only advise against, a SHOULD or RECOMMENDED, and
rejecting those would refuse documents that are legal. Notion's
specification, in `examples/notion/openapi-undocumented.json`, has
`"enum": []` twice: JSON Schema 2020-12 says an enum's array "SHOULD have at
least one element", and an empty one allows no value at all. Downstream it
was taken for a plain string and merged with five unrelated ones.

A `ValidateStrict` beside `Validate` would also report what is only advised
against: an empty `enum`, an `example` its schema would reject for more than its
kind, which `Validate` already checks, a
response or parameter without a `description`. `Validate` stays as it is,
so a document that loads today still validates.

Open questions:

- whether strict validation returns every finding or stops at the first, as
  `Validate` does;
- which advice to check first; the empty `enum` has a real case.
