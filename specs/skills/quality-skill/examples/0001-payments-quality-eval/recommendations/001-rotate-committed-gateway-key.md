---
type: Recommendation
title: Rotate and remove the committed gateway credential
description: Reference Recommendation — close the committed live-credential gap holding Sparrow Payments at Unacceptable.
tags: [skill, quality, evaluation, example, security]
timestamp: 2026-06-17T00:00:00Z
---

> **Reference instance — non-normative.** A captured example of a single
> recommendation artifact the [`/quality` skill](../../../quality-skill.md)
> emits alongside its [report](../report.md). It is written to stand on its own:
> a reader can triage and route it without the report or the session.

# Rotate and remove the committed gateway credential

**Target / factor:** Sparrow Payments API (root) → Security → Secrets handling
**In-scope requirement:** *No credentials are committed to the repository*
**Current rating:** Unacceptable — **binding constraint** on the whole-model rating.

## Gap

A live payment-gateway credential is committed to the repository in plaintext.
A committed working secret is an immediate exposure, which meets the Unacceptable
criterion for the Secrets handling sub-factor and holds the entire model at
Unacceptable.

**Evidence** (credential referenced by location and type only; the value is not
reproduced here):

- `internal/gateway/client.go:48` — a live **gateway secret key** assigned to a
  package-level constant. Matches the active-key format; not a placeholder.
- Present in committed history, so removal from `HEAD` alone does not contain the
  exposure — the key must be treated as compromised and rotated.
- Not in scope of this gap: `internal/gateway/client_test.go:12` holds a **test
  publishable key**, which is non-secret by design.

## Options

- **(a) Rotate at the gateway, then move the secret to runtime configuration.**
  Issue a new key in the gateway console, revoke the exposed one, and load the
  secret from the environment / secret manager at startup instead of source.
- **(b) Move the secret to configuration without rotating.** Removes it from the
  working tree but leaves the exposed key valid in history — does not contain the
  exposure.
- **(c) Purge the key from git history (e.g. history rewrite) without rotating.**
  Reduces discoverability but cannot guarantee the key was never copied; the
  exposed key stays valid.

## Recommended

**(a) Rotate at the gateway, then move the secret to runtime configuration.**
Rotation is required, not merely removal: the key is already exposed in history,
so only revoking it contains the risk. Options (b) and (c) leave a valid exposed
credential and do not clear the Unacceptable criterion.

## Done-criterion

The requirement *No credentials are committed to the repository* reaches
**Target** against its criterion — no live credential is present in the working
tree, and the previously exposed key has been revoked at the gateway so it can no
longer be used. A later `improve` re-evaluates this scope to confirm the rating
moved off Unacceptable; clearing this constraint is expected to lift root
Security to **Target** and the whole-model rating off the floor to **Minimum**,
where the webhook-delivery deduplication gap (see
[003](003-bound-webhook-dedup-window.md)) becomes the binding constraint until it
too is closed.
