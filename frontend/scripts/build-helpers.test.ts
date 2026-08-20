import { describe, expect, test } from "bun:test";
import { rewriteIndexHTMLContents } from "./build-helpers";

describe("rewriteIndexHTMLContents", () => {
  test("rewrites the development script path to the built entry", () => {
    const original = `<html><head></head><body><script type="module" src="../src/index.tsx"></script></body></html>`;
    const updated = rewriteIndexHTMLContents(original);
    expect(updated).toContain("/index.js");
    expect(updated).not.toContain("../src/index.tsx");
    expect(updated).toContain('href="/index.css"');
  });

  test("appends a stylesheet when the document has no head close tag", () => {
    const original = `<script type="module" src="../src/index.tsx"></script>`;
    const updated = rewriteIndexHTMLContents(original);
    expect(updated).toContain("/index.js");
    expect(updated).toContain('href="/index.css"');
  });
});
