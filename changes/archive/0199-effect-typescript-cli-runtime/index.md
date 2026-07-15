# Effect TypeScript CLI runtime

Children of the
[Effect TypeScript CLI runtime](../0199-effect-typescript-cli-runtime.md)
Change Case.

Implementation is complete and ready for review. The runtime uses Bun 1.3.14,
Effect v4 beta.98 from the reviewed Effect `main` branch, and distinct glibc and
musl Linux release assets.

# Concepts

- [Functional spec](spec.md) - requirements for the clean Go-to-TypeScript
  cutover, CLI and artifact compatibility, graceful source resolution,
  isolated agent contexts, evaluator capabilities, distribution continuity,
  and acceptance gates.
- [Design doc](design.md) - the Effect v4 architecture, pure domain core,
  provider SDK adapters, context lifecycle, resolved runtime-spike decisions,
  standalone executable packaging, release cutover, and verification strategy.
