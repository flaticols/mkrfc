---
title: "{{ .Title }}"
authors:
  - "{{ .Author }}"
reviewers: []
status: "DRAFT"
date: "{{ .Date }}"
{{ if .PR -}}
pr: "{{ .PR }}"
{{ end -}}
tags: []
---

# [RFC] {{ .Title }}

**Authors:**

- {{ .Author }}

| Field  | Value       |
| ------ | ----------- |
| Status | **DRAFT**   |
| Date   | {{ .Date }} |
{{ if .PR -}}
| PR     | [View PR]({{ .PR }}) |
{{ end -}}
{{ if .TargetServices -}}
| Target Services | {{ .TargetServices }} |
{{ end -}}

---

## Executive Summary

{{ .Summary }}
{{ if index .Sections "motivation" }}
---

## Motivation

<!-- TODO: why is this change needed? -->
{{ end -}}
{{ if index .Sections "implementation" }}
---

## Proposed Implementation

<!-- TODO: describe your proposed implementation -->
{{ end -}}
{{ if index .Sections "metrics" }}
---

## Metrics & Dashboards

<!-- TODO: define key metrics and how to measure them -->
{{ end -}}
{{ if index .Sections "drawbacks" }}
---

## Drawbacks

<!-- TODO: list reasons why we should not do this -->
{{ end -}}
{{ if index .Sections "alternatives" }}
---

## Alternatives

<!-- TODO: describe other ways to achieve the same outcome -->
{{ end -}}
{{ if index .Sections "impact" }}
---

## Potential Impact and Dependencies

<!-- TODO: identify affected teams, systems, and dependencies -->
{{ end -}}
{{ if index .Sections "unresolved" }}
---

## Unresolved Questions

<!-- TODO: list open questions and areas still being defined -->
{{ end -}}
{{ if index .Sections "conclusion" }}
---

## Conclusion

<!-- TODO: summarize why this is the right decision to make now -->
{{ end -}}
