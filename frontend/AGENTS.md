# Flower frontend

The Flower frontend is a SolidJS application that fronts the Flower Go API.

## Tech

* SolidJS
* Tailwind v4
* TypeScript
* Bun

There is no Vite and no Node. Use Bun for install, test, bundle, and the dev server.

## Design guide

See [docs/reference/frontend-design-guide.md](../docs/reference/frontend-design-guide.md).

## Testing

* Use Bun's test runner with `@solidjs/testing-library` and happy-dom.
* Do not talk to live services from unit tests. Mock the API boundary.
* Tests must exercise the same production components as runtime code.
