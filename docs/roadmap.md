# Roadmap

Work not yet done. An entry is deleted once it is.

## More strict checks

`Document.ValidateStrict` reports what the specifications allow but a document
rarely means: an unknown format, and a key outside the components that would
not do as a component's. Other advice it could check, each with a real case
first:

- an empty `enum`, which JSON Schema 2020-12 says "SHOULD have at least one
  element" and which allows no value at all; Notion's undocumented
  specification had it twice, and downstream it was taken for a plain string;
- an `example` its schema would reject for more than its kind, which
  `Validate` already checks;
- a response or parameter without a `description`.
