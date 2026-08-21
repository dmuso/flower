/** @jsxImportSource solid-js */
import { afterEach, beforeEach, describe, expect, it, mock } from "bun:test";
import { cleanup, fireEvent, render, screen } from "@solidjs/testing-library";
import { Board } from "./Board";

beforeEach(() => cleanup());
afterEach(() => cleanup());

function renderBoard() {
  return render(() => (
    <Board
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
}

describe("Board", () => {
  it("renders four empty columns with empty-state copy and 0 / 10", () => {
    renderBoard();
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
    expect(screen.getByRole("button", { name: "Keyboard shortcuts" })).toBeTruthy();
    expect(screen.queryByRole("button", { name: "Add a story", hidden: false })).toBeTruthy();
    const backlog = screen.getByTestId("column-backlog");
    expect(backlog.textContent || "").not.toContain("Add a story");
  });

  it("opens a Close-only shortcut overlay from ? and binds A to Add a story", () => {
    const click = mock(() => undefined);
    renderBoard();
    const add = screen.getByRole("button", { name: "Add a story" });
    add.addEventListener("click", click);
    fireEvent.keyDown(screen.getByTestId("board"), { key: "A" });
    expect(click).toHaveBeenCalled();
    fireEvent.click(screen.getByRole("button", { name: "Keyboard shortcuts" }));
    expect(screen.getByRole("dialog", { name: "Keyboard shortcuts" })).toBeTruthy();
    expect(screen.getByRole("button", { name: "Close" })).toBeTruthy();
    fireEvent.click(screen.getByRole("button", { name: "Close" }));
    expect(screen.queryByRole("dialog", { name: "Keyboard shortcuts" })).toBeNull();
  });
});
