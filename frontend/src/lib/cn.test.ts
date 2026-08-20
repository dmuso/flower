import { describe, expect, it } from "bun:test";
import { cn } from "./cn";

describe("cn", () => {
  it("merges conflicting Tailwind classes", () => {
    expect(cn("px-4", "px-2")).toBe("px-2");
  });

  it("drops falsy conditional classes", () => {
    const hidden = false;
    expect(cn("rounded-lg", hidden && "hidden", "border")).toBe("rounded-lg border");
  });
});
