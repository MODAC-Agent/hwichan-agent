# Week 4 Improvement Plan for gosu-review

This document summarizes the refinement strategy for the `gosu-review` skill, based on the real-world code patterns identified in the `skipjd` test data.

## 1. Context-Aware Rule Refinement

### Error Wrapping Strategy (Rule 48)
- **Observation:** In `crawler.go` and `output.go`, internal I/O or parsing errors were wrapped with `%w`.
- **Improvement:** Update Rule 48 to distinguish between "domain errors" (caller needs to branch) and "internal implementation errors".
- **Action:** Add a check to see if the error is likely to be handled by the caller (e.g., returned from a public method) vs. a terminal/internal failure (e.g., `fmt.Fprintf`).

### Guard Clause Logic (Rule 2)
- **Observation:** `detail_scraper.go` contains nested `switch` and `if` blocks.
- **Improvement:** Refine Rule 2 to handle `switch` cases. Sometimes a `switch` is cleaner than a long list of `if` guard clauses.
- **Action:** Define a threshold for "meaningful nesting" vs. "unnecessary complexity" in `switch` statements.

### Receiver Type Selection (Rule 42)
- **Observation:** The `Crawler` struct uses pointer receivers, but some methods might not mutate state.
- **Improvement:** Ensure the rule checks for **consistency** across the entire struct rather than just flagging individual methods.
- **Action:** Enforce the "all pointers or all values" rule more strictly when reviewing a single file.

## 2. Eliminating Outdated/Noisy Advice

### Range Iteration & Go 1.22+ (Rule 31)
- **Observation:** The skill gave noisy advice about `range` value copies even for small structs or cases where Go 1.22+ optimizations apply.
- **Improvement:** 
    - Add a **size threshold** for flagging value copies (e.g., ignore if struct is small).
    - Add a note about **Go 1.22+ loop variable semantics** to prevent outdated advice.
- **Action:** Update `references/trigger-07-range-iteration.md` with explicit exclusion criteria.

## 3. Handling Designing Decisions (Rule 7/11)

### Functional Options Pattern (Rule 11)
- **Observation:** `crawler.go` implements Functional Options but the skill might still suggest it if the implementation is non-idiomatic.
- **Improvement:** Recognize existing patterns to avoid redundant suggestions.
- **Action:** Add a "Pattern Detection" step before suggesting Rule 11.

## 4. Safety & PII Awareness
- **Observation:** Real code often contains "Key" or "Token" as technical terms.
- **Improvement:** Reduce false positives in "Security" checks by recognizing common Go/GORM patterns (e.g., `SourceKey`).
- **Action:** Whitelist common technical suffixes/patterns in the skill's initial scan.

---

### Target Files for Testing
Use these files in `week3/testdata/` to verify the refinements:
- `crawler/crawler.go`: Dependency injection, Error wrapping.
- `gamejob/detail_scraper.go`: Guard clauses, `defer` evaluation.
- `model/job_posting.go`: Struct design, `range` performance.
