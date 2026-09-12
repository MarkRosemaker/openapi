The primary goals of this package are:

- **Parsing** OpenAPI specifications into a structured format.
- **Formatting** the parsed specifications, including sorting maps and merging duplicate content.
- **Adding information programmatically** to the specifications.
- **Marshalling** the modified specifications back into their original format.
- **Utilizing** the parsed specification for code generation.

This module is deliberately kept focused on representing and validating a
specification. Transformations that not everyone needs — flattening, deduplicating,
enriching, generating code — live in [separate modules](#the-openapi-family) so that
users who only want to parse, validate, and prettify a spec don't pay for them.
