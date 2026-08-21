/** @jsxImportSource solid-js */
import { afterEach, beforeEach, describe, expect, it, mock } from "bun:test";
import { Router, Route } from "@solidjs/router";
import { cleanup, fireEvent, render, screen } from "@solidjs/testing-library";
import { SignInPage } from "./SignIn";

const originalFetch = globalThis.fetch;

beforeEach(() => cleanup());
afterEach(() => {
  cleanup();
  globalThis.fetch = originalFetch;
});

describe("SignInPage", () => {
  it("renders sign in and magic-link secondary", () => {
    render(() => (
      <Router>
        <Route path="/" component={SignInPage} />
      </Router>
    ));
    expect(screen.getByRole("heading", { name: "Sign in" })).toBeTruthy();
    expect(screen.getByRole("button", { name: "Sign in" })).toBeTruthy();
    expect(screen.getByRole("button", { name: "Email me a link instead" })).toBeTruthy();
  });

  it("shows mismatch copy", async () => {
    globalThis.fetch = mock(async () => {
      return new Response(JSON.stringify({ error: { code: "unauthorized", message: "That email and password don’t match." } }), {
        status: 401,
        headers: { "content-type": "application/json" },
      });
    }) as unknown as typeof fetch;
    render(() => (
      <Router>
        <Route path="/" component={SignInPage} />
      </Router>
    ));
    fireEvent.input(screen.getByLabelText("Email"), { target: { value: "maya@example.com" } });
    fireEvent.input(screen.getByLabelText("Password"), { target: { value: "nope" } });
    fireEvent.click(screen.getByRole("button", { name: "Sign in" }));
    expect((await screen.findByTestId("signin-error")).textContent).toBe("That email and password don’t match.");
  });
});
