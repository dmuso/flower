import { describe, expect, it } from "bun:test";
import { readFileSync } from "node:fs";
import { join } from "node:path";

describe("committed tailwind.css", () => {
  it("includes board and primary utilities so class-name tests cannot go green alone", () => {
    const css = readFileSync(join(import.meta.dir, "tailwind.css"), "utf8");
    expect(css).toContain(".bg-bloom");
    expect(css).toContain(".min-w-\\[22rem\\]");
  });
});
