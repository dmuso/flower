/** @jsxImportSource solid-js */
import { afterEach, beforeEach, describe, expect, it } from "bun:test";
import { Router, Route } from "@solidjs/router";
import { cleanup, render, screen } from "@solidjs/testing-library";
import { NameProjectPage } from "./NameProject";

beforeEach(() => cleanup());
afterEach(() => cleanup());

describe("NameProjectPage", () => {
  it("asks for the first project name", () => {
    render(() => (
      <Router>
        <Route path="/" component={NameProjectPage} />
      </Router>
    ));
    expect(screen.getByText("Name the first project.")).toBeTruthy();
    expect(screen.getByLabelText("Project name")).toBeTruthy();
    expect(screen.getByRole("button", { name: "Create project" })).toBeTruthy();
  });
});
