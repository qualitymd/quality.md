# QUALITY.md Format

QUALITY.md is a plain text represenation of a quality model. It can be used to specify and evalute the quality requirements for a software system or component.

A QUALITY.md file contains two parts: YAML frontmatter with the structured quality model and the markdown body.

## Quality Model
The quality model is embedded in tye YAML front matter at the begninnign of the file. The front matter block must begin with a line containing exactly --- and end with a line containing exactly ---. The YAML content between these delimiters is parsed according to the schema defined below.

Example:
```yaml

factors:
  security:
    requirements:
      "meets app sec security standars":
        rules: "./standards/appsec-checklist.md"
        rating:
          pass: "satisfies all checks"
  maintainability:
    factors:
      reusability:
        requirements:
          "uses types from common package":
            rules: >
              TODO: inline rule a
              TODO: inline rule b
            rating:
              pass: "inferentiaal pass criteria"
              fail: "inferentiaal fail criteria"
      testability:
        requirements:
          "testability requirement A":
            bash: "pnpm test:unit"

```

### Schema
