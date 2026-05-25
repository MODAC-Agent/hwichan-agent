# hwichan-agent

A 5-week study log of building **[skipJD](https://github.com/oak-cassia/skipjd)** — a job posting recommendation system written in Go, with an emphasis on idiomatic Go practices. The system itself lives in a separate repository (`skipjd`).

## Goals of the Study

- Move past "code that works" toward **idiomatic Go code**. Syntactically clean output can still hide anti-patterns; reducing those was the starting point.
- Use the project itself as a forcing function for studying *when* an LLM belongs in a system and when it does not — recording each apply / abandon decision as the study's main output.

## What I Built — skipJD

A recommendation system targeting the gamejob portal. A crawler collects postings; an OCR worker converts image-only bodies into text; an extractor turns each body into a structured `experience / competency / trait` JSON; the matcher compares the result against user profiles.

```
[gamejob crawler]
  │  HTML body  → job_posting_bodies
  │    · no image      → ready_for_llm = true
  │    · has image     → ready_for_llm = false
  │  image URLs → job_posting_images
  ▼
[OCR worker (LLM)]
  Targets postings that have images but no OCR body yet.
  Stores HTML body + "[OCR]" marker + OCR text as a merged body,
  then flips ready_for_llm = true.
  ▼
job_posting_bodies  (only rows with ready_for_llm = true proceed)
  ▼
[extractor (LLM)]
  body → experience / competency / trait JSON
  ▼
[structured matching input]
  + experience / competency / trait
  + job codes (multi-valued)
  + normalized company ID
  + user preferences · years of experience
  ▼
[recommendation matcher] → [user recommendations]
```

The three batch workers are gated by the `ready_for_llm` flag on each body row, so they do not depend on execution order or concurrency.

Repository layout (skipJD):

- `cmd/` — batch and service entry points (`crawler`, `ocr-worker`, `extractor`, `pipeline`, `api`, `seed-preferences`, `normalize-companies`, `user-extractor`, `notify`)
- `internal/` — domain packages (`crawler`, `gamejob`, `ocrworker`, `extractor`, `matcher`, `repository`, `model`, `geminiexec`, `retry`, `batch`, `mailing`, `telemetry`, …)

## Final Retrospective

[*"LLM을 잘 활용하려다 다시 마주친 문제 정의의 중요성"* (Medium)](https://medium.com/@gnlckswjd1/llm%EC%9D%84-%EC%9E%98-%ED%99%9C%EC%9A%A9%ED%95%98%EB%A0%A4%EB%8B%A4-%EB%8B%A4%EC%8B%9C-%EB%A7%88%EC%A3%BC%EC%B9%9C-%EB%AC%B8%EC%A0%9C-%EC%A0%95%EC%9D%98%EC%9D%98-%EC%A4%91%EC%9A%94%EC%84%B1-750b721312e4) — wrap-up of the 5-week study.

## What Went Wrong (and What I Learned)

Two LLM ideas I designed and then **abandoned**, plus one piece of work I deferred two weeks in a row — the *judgment process* behind each of these is the actual output of this study, more than the artifacts themselves.

1. **A general-purpose AI agent that crawls any URL** — My first take blamed token cost. The honest reason came later: once I committed to targeting a single recruiting portal, the page structures I had to handle shrank to a small, finite set that code could cover, and **the LLM's adaptability stopped producing value**.
2. **An LLM pipeline to normalize user-entered company names** — Designed end-to-end (dictionary/cache → similarity Top-N → LLM precision matching) before scrapping it. If the input form lets users pick from the company list the crawler already maintains, **the problem itself disappears**.
3. **Deferring the recommendation matcher two weeks in a row** — The first deferral went into shaping the matching inputs; the second went into making the extraction outputs reliable. Both came from the same call: *stabilize the data layer first*.

## The LLM-Application Criterion I Settled On

My default heuristic was *"unstructured data → LLM, structured data → code"*. Building skipJD revealed there's a question that has to come **before** that one:

> **Can the set of deterministic problems be narrowed enough that code can handle them all, or is it so open-ended that it's effectively a non-deterministic problem?**

- The set can be narrowed to something code can handle → **use code**. Cost, reproducibility, and verifiability all favor code.
- The set is effectively infinite, or finite but too large → two options: (a) hand it off to an LLM, or (b) **narrow it back into code's domain by reshaping the problem**. "Closing the input space" in the company-name case and "single-portal targeting" in the crawler case are both (b).
- So the question that comes before *"should I use an LLM?"* is **"can I narrow the problem?"** If yes, the LLM step disappears entirely.

In skipJD, the LLM is left in exactly two places — extractor (job body → structured JSON) and OCR (image → text). Both are *unstructured → structured* input-stage transformations, and both call the same tool (`gemini-cli`). Company-name normalization, matching keys, and retry/timeout policies all ended up in code.
