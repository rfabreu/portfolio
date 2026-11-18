# Contact Notification Labels — Design Spec

> Clean field names in the Netlify Forms email notification for the `contact` form, with no visible change on the site.

---

## 1. Overview & Objectives

The owner reports the notification is badly formatted. Netlify builds the names it shows from the form's `<label>` text, and the labels were `NAME *`, `EMAIL *`, `PROJECT TYPE *`, `TARGET TIMELINE`, `PROJECT BRIEF *`: raw uppercase text with a literal ` *` inside the label.

Goals:
1. Label text is clean sentence case (`Name`, `Email`, `Project type`, `Target timeline`, `Project brief`) with no asterisk.
2. The page looks the same: uppercase via CSS, the asterisk via a CSS pseudo-element.
3. Field `name` attributes, `required` attributes, form wiring and the honeypot are untouched.

Non-goals: copy changes, new fields, restyling, changing `name` attributes (Netlify stores submissions by name; renaming would split existing data), P10 positioning changes (no budget or pricing fields).

---

## 2. Research: how Netlify derives labels

Sources read on 2026-11-18:

- Setup docs, https://docs.netlify.com/manage/forms/setup.md: form detection scans deploy HTML for `data-netlify="true"` (or `netlify`), the `name` attribute names the form, the build strips `data-netlify` and injects a hidden `form-name` input; JS-rendered forms are not detected. Says **nothing** on label derivation.
- Notifications docs, https://docs.netlify.com/forms/notifications/: only recommends `name="email"` for Reply-to. Says nothing on labels.
- Netlify Support Forums (community and staff replies, **not** official docs):
  - https://answers.netlify.com/t/how-does-netlify-forms-determine-labels-in-email-notifications/9255 : staff say `for` attributes are not handled for certain fields and recommend wrapping the input in the label; duplicate labels clobber each other.
  - https://answers.netlify.com/t/netlify-forms-do-not-respect-form-labels-in-emails/41802 : a poster's wrapped label `Full name:` came out as `Name:` (matches the `name` attribute), but no reply states the rule.
  - https://answers.netlify.com/t/form-is-being-displayed-differently-on-netlify-and-emails-than-it-is-inside-payload-and-in-html/76945 : wrapping inputs in labels fixed mismatched labels.

**The docs do not settle it.** Conclusion: label text is the likely source, with the `name` attribute as a possible alternative, and nothing official says how inner elements or whitespace in a label are handled. The owner's report (asterisks in the notification) is consistent with label text being used.

Approach correct under both readings:
- If label text is used: it is now clean.
- If the `name` attribute is used: we change nothing there, and renaming would break stored data, so that reading is out of scope.
- Because inner-element handling is undocumented, the asterisk is **not** moved into a child element (even `aria-hidden`, which still counts as text content); a CSS `::after` pseudo-element is not part of the DOM text at all, so it works however Netlify parses labels.

---

## 3. Design

| Field (`name`) | Old label text | New label text | Required marker |
|---|---|---|---|
| `name` | `NAME *` | `Name` | yes |
| `email` | `EMAIL *` | `Email` | yes |
| `project_type` | `PROJECT TYPE *` | `Project type` | yes |
| `timeline` | `TARGET TIMELINE` | `Target timeline` | none (optional field) |
| `message` | `PROJECT BRIEF *` | `Project brief` | yes |

- Each `<label>` gains Tailwind `uppercase` (renders the same capitals as before) and, for required fields, a `field-required` class.
- `global.css` adds `.field-required::after { content: " *"; }`, matching the old text ` *` (space, asterisk) and inheriting the label colour in both themes.
- Visible labels stay exactly `NAME`, `EMAIL`, `PROJECT TYPE`, `TARGET TIMELINE`, `PROJECT BRIEF` (the task's example names "Timeline" and "Message" would have changed visible copy, so the existing wording is kept in sentence case).
- Unchanged: `for`/`id` pairs, every `name`, every `required`, the honeypot `<label>Don't fill this out: <input name="bot-field" /></label>` (hidden, Netlify's documented pattern), `form-name`.
- Accessibility: `required` on the controls is what announces required state; the old literal `*` was text, the new one is generated content, both may be read as "star" by some readers. No `aria-required` added, not needed with native `required`.

---

## 4. Verification

Vitest `tests/contact-labels.test.ts`: label text per field, no `*`, no child elements, `name` and `required` unchanged, marker class only on required fields, `uppercase` present, honeypot untouched, CSS rule present. Build checks of `dist/index.html`. Screenshots of the form before (origin/main build) and after, light and dark, 1280px and 390px: pixel diff is 37 to 39 pixels (max channel delta 24 of 255) confined to one 8x8 box at the first label's asterisk, i.e. sub-pixel anti-aliasing, nothing else differs.

---

## 5. Open questions

1. The real notification can only be confirmed by the owner submitting the form after deploy. Netlify's label derivation is undocumented, so this fix is a best-supported inference, not proven.
2. If labels still look wrong, the fallback is wrapping each input inside its `<label>` (staff recommendation in thread 9255). Not done here because it changes the DOM and spacing and the owner's report shows the `for` labels are being read.
3. Notification field order and the `project_type` slug are Netlify's; not changeable without renaming fields.
