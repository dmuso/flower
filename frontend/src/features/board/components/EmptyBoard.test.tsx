/** @jsxImportSource solid-js */
import { afterEach, beforeEach, describe, expect, it } from "bun:test";
import { cleanup, render, screen } from "@solidjs/testing-library";
import { EmptyBoard } from "./EmptyBoard";

beforeEach(() => cleanup());
afterEach(() => cleanup());

describe("EmptyBoard", () => {
  it("renders four empty columns with slice 0 copy and 0 / 10", () => {
    render(() => (
      <EmptyBoard
        organisation="Acme"
        project="Trail"
        pack={{
          current_points: 0,
          denominator: 10,
          velocity_source: "initial",
          current_window_ends_at: "2026-08-23T14:00:00Z",
        }}
        onSignOut={() => undefined}
      />
    ));
    expect(screen.getByText("Acme")).toBeTruthy();
    expect(screen.getByText("Trail")).toBeTruthy();
    expect(screen.getByText("Nothing waiting. Capture a story.")).toBeTruthy();
    expect(screen.getByText("Nothing ranked yet. Pull from Icebox when you’re ready.")).toBeTruthy();
    expect(screen.getByText("Nothing in this window yet. Pull a story from Icebox or create one.")).toBeTruthy();
    expect(screen.getByText("Nothing in Done yet. Accepted work lands here after it ages out of this window.")).toBeTruthy();
    expect(screen.getByRole("button", { name: "Add a story" })).toBeTruthy();
    expect(screen.getByText(/0 \/ 10/)).toBeTruthy();
    expect(document.body.innerHTML).toContain("bg-paper");
    expect(document.body.innerHTML).toContain("bg-bloom");
    expect(document.body.innerHTML).not.toMatch(/bg-(indigo|purple|violet|cyan)-/);
  });
});
