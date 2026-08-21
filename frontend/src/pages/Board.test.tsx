/** @jsxImportSource solid-js */
import { describe, expect, it } from "bun:test";
import { BoardPage } from "./Board";

describe("BoardPage", () => {
  it("exports the board page", () => {
    expect(typeof BoardPage).toBe("function");
  });
});
