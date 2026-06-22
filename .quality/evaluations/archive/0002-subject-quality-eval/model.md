# Active model (snapshot)

- **Altitude:** subject
- **Active model:** the repository `QUALITY.md` (snapshot taken at run start).
- **Subject:** the QUALITY.md project deliverables it targets.
- **Scope:** whole model — both targets, every in-scope requirement.
- **Effort:** standard.
- **Rating scale:** outstanding > target > minimum > unacceptable (best to worst).

## Targets and requirements assessed

### Target `format-spec` (source `./SPECIFICATION.md`)

- Direct requirement: *the format specification is complete*
- Factor **clarity**
  - *the format specification admits a single interpretation*
  - *the format specification separates rules from rationale*
  - *the format specification defines its terms before use*
- Factor **consistency**
  - *the format specification is internally consistent*
- Factor **verifiability**
  - *each rule is observable or testable*
  - *the format's constructs are shown with valid and invalid examples*
- Factor **extensibility**
  - *the format specifies its core and how it extends and evolves*
- Factor **usability**
  - *the format specification is well-structured and readable*

### Target `readme` (source `./README.md`)

- Factor **approachability**
  - *the README says what QUALITY.md is and who it's for*
  - *the README shows the format and its payoff by example*
  - *the README gets a newcomer to a first result quickly*
  - *the README reflects what the CLI and spec actually provide*

The model root declares no own requirements or factors; it is a grouping node
whose aggregate considers only its two child targets.
