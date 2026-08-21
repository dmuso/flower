/** @jsxImportSource solid-js */
import { afterEach, beforeEach, describe, expect, it } from "bun:test";
import { Router, Route } from "@solidjs/router";
import { cleanup, render, screen } from "@solidjs/testing-library";
import { CheckEmailPage } from "./CheckEmail";

beforeEach(() => cleanup());
afterEach(() => cleanup());

describe("CheckEmailPage", () => {
  it("shows password-signup verify copy", () => {
    render(() => (
      <Router>
        <Route path="/" component={CheckEmailPage} />
      </Router>
    ));
    expect(screen.getByText("Check your email to verify.")).toBeTruthy();
  });
});
