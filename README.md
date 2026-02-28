# mkrfc

RFC management CLI for engineering teams.

```
go install github.com/flaticols/mkrfc@latest
```

## Usage

```
mkrfc              New RFC if none exist, otherwise list all RFCs
mkrfc new          Create a new RFC with interactive form
mkrfc list         List all RFCs
mkrfc resolve      Transition a draft RFC (accept / reject / supersede)
mkrfc review       Add a reviewer; optionally approve or reject
mkrfc templates    List available templates
mkrfc version      Print version
mkrfc llm          Print compact LLM-friendly help in XML format
mkrfc help         Show help
```

## How it works

RFCs are stored as Markdown files with YAML frontmatter in `docs/rfc/`:

```
docs/rfc/
  0001-add-async-job-queue.md
  0002-migrate-to-postgresql.md
  .templates/          ← custom templates (optional)
    lightweight.md
```

### `mkrfc new`

Interactive TUI form:
- **Template** — picker shown when multiple templates exist
- **RFC Title** — required
- **Summary** — brief description
- **Optional Sections** — multi-select: Motivation, Proposed Implementation, Drawbacks, Alternatives, Metrics, Impact, Unresolved Questions, Conclusion

Auto-detected from the environment:
- **Author** — from `git config user.name` + `git config user.email`
- **Date** — today's date
- **Number** — next sequential 4-digit number
- **PR** — current branch's pull request URL (via `gh pr view`, if available)

### `mkrfc list`

```
mkrfc list
mkrfc list --status DRAFT
```

### `mkrfc resolve`

Picks a draft RFC, selects a new status (ACCEPTED / REJECTED / SUPERSEDED), and prompts for a resolution note. Appends a `## Resolution` section and updates the frontmatter.

### `mkrfc review`

Picks a draft RFC, pre-fills your name from git config, and lets you:
- Add yourself as a reviewer
- Approve (sets ACCEPTED + resolution note)
- Reject (sets REJECTED + resolution note)

`mkrfc approve` is an alias for `mkrfc review`.

## Templates

The default template is embedded in the binary. Add custom templates by dropping `.md` files (Go `text/template` syntax) into `docs/rfc/.templates/`. When multiple templates exist, `mkrfc new` shows a picker.

Template data fields: `Title`, `Author`, `Date`, `Summary`, `PR`, `Sections`, `TargetServices`, `Motivation`, `Implementation`, `Metrics`, `Drawbacks`, `Alternatives`, `Impact`, `Unresolved`, `Conclusion`.

## Frontmatter schema

```yaml
title: "RFC Title"
authors:
  - "Name <email>"
reviewers: []
status: DRAFT          # DRAFT | ACCEPTED | REJECTED | SUPERSEDED
date: "2024-03-10"
resolved_date: ""      # set by resolve / review
pr: ""                 # PR URL, auto-detected from gh CLI
tags: []
```

## License

[MIT](LICENSE)
