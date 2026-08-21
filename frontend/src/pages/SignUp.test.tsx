/** @jsxImportSource solid-js */
import { afterEach, beforeEach, describe, expect, it, mock } from "bun:test";
import { Router, Route } from "@solidjs/router";
import { cleanup, fireEvent, render, screen } from "@solidjs/testing-library";
import { SignUpPage } from "./SignUp";

const originalFetch = globalThis.fetch;

beforeEach(() => {
  cleanup();
});

afterEach(() => {
  cleanup();
  globalThis.fetch = originalFetch;
});

function renderSignUp() {
  return render(() => (
    <Router>
      <Route path="/" component={SignUpPage} />
    </Router>
  ));
}

describe("SignUpPage", () => {
  it("asks for email and password and does not ask for a username", () => {
    renderSignUp();
    expect(screen.getByText("Email and a password. That’s enough to start.")).toBeTruthy();
    expect(screen.getByRole("button", { name: "Sign up" })).toBeTruthy();
    expect(screen.getByRole("button", { name: "Email me a link instead" })).toBeTruthy();
    expect(screen.getByText("Already a user?")).toBeTruthy();
    expect(screen.queryByLabelText(/username/i)).toBeNull();
    expect(screen.getByLabelText("Email")).toBeTruthy();
    expect(screen.getByLabelText("Password")).toBeTruthy();
  });

  it("flips to magic link and shows Email is enough.", () => {
    renderSignUp();
    fireEvent.click(screen.getByRole("button", { name: "Email me a link instead" }));
    expect(screen.getByRole("button", { name: "Email me a link" })).toBeTruthy();
    expect(screen.getByRole("button", { name: "Use a password instead" })).toBeTruthy();
    expect(screen.getByText("Email is enough.")).toBeTruthy();
    expect(screen.queryByText("Email and a password. That’s enough to start.")).toBeNull();
    expect(screen.queryByLabelText("Password")).toBeNull();
  });

  it("shows taken-email copy with Sign in as the next action", async () => {
    globalThis.fetch = mock(async () => {
      return new Response(JSON.stringify({ error: { code: "email_taken", message: "That email already belongs to a user." } }), {
        status: 409,
        headers: { "content-type": "application/json" },
      });
    }) as unknown as typeof fetch;
    renderSignUp();
    fireEvent.input(screen.getByLabelText("Email"), { target: { value: "maya@example.com" } });
    fireEvent.input(screen.getByLabelText("Password"), { target: { value: "secret12" } });
    fireEvent.click(screen.getByRole("button", { name: "Sign up" }));
    expect((await screen.findByTestId("signup-error")).textContent).toBe("That email already belongs to a user.");
    expect(screen.getByRole("link", { name: "Sign in" })).toBeTruthy();
    expect(screen.queryByRole("button", { name: "Sign up" })).toBeNull();
  });

  it("uses only the locked palette classes", () => {
    renderSignUp();
    const html = document.body.innerHTML;
    expect(html).toContain("bg-bloom");
    expect(html.includes("paper") || html.includes("bg-paper")).toBe(true);
    expect(html).not.toMatch(/bg-(indigo|purple|violet|cyan|teal|orange|lime)-/);
  });
});
