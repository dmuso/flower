/** @jsxImportSource solid-js */
import { afterEach, beforeEach, describe, expect, it } from "bun:test";
import { Router, Route } from "@solidjs/router";
import { cleanup, render, screen } from "@solidjs/testing-library";
import { HomePage } from "./Home";

beforeEach(() => {
  cleanup();
});

afterEach(() => {
  cleanup();
});

describe("HomePage", () => {
  it("renders sign up copy without a username field", () => {
    render(() => (
      <Router>
        <Route path="/" component={HomePage} />
      </Router>
    ));
    expect(screen.getByRole("heading", { name: "Sign up" })).toBeTruthy();
    expect(screen.getByText("Email and a password. That’s enough to start.")).toBeTruthy();
    expect(screen.queryByLabelText(/username/i)).toBeNull();
  });
});
