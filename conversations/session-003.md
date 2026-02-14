# Session 3 — 2026-02-14

## Summary

Established learning interaction rules and Claude Code instructions.

## Changes

- Created `CLAUDE.md` with interaction rules for the learning project
- Defined workflow: Denis announces task → Claude explains → Denis implements → Claude writes tests → Denis iterates → Claude reviews

## Discussion

- Denis raised the concern that if Claude does all the work, he won't learn
- Agreed on guided implementation as the default mode
- Claude writes tests based on requirements (not implementation), Denis writes production code
- Escalation path when stuck: hint → partial example → full solution

---

## Full Transcript

**Denis:**
> Let's talk about you and our interactions. I want to learn about go, rust and kubernetes. Do you have suggestions on how we can interact so that I can learn ? If you do all the work I don't think this gonna work for me

**Claude:**
> Suggested several interaction modes:
> 1. **Guided implementation** — Claude explains, Denis writes, Claude reviews
> 2. **Challenge mode** — Claude writes failing test, Denis makes it pass
> 3. **Explain-then-implement** — Claude explains new concepts, Denis implements exercises
> 4. **Code review mode** — Denis writes, Claude reviews for idiomatic patterns
>
> Recommended: Denis writes the code, Claude guides and reviews. Claude only writes boilerplate/config. Escalation: hint → partial example → full solution.

**Denis:**
> yes save them to CLAUDE.md

**Claude:**
> Created `CLAUDE.md` with all interaction rules.

**Denis:**
> I'd like you to write the tests too, but not based on my implementation, rather based on the requirements. Every time I start working on something, I'll tell you what I'm about to do. Give me the explanation for new concepts at first and let me implement it, you'll write the test after

**Claude:**
> Updated `CLAUDE.md` workflow:
> 1. Denis announces what he's working on
> 2. Claude explains new concepts
> 3. Denis implements
> 4. Claude writes tests based on requirements
> 5. Denis iterates until tests pass
> 6. Claude reviews for idiomatic patterns

**Denis:**
> commit and push
