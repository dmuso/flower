/** @jsxImportSource solid-js */
import { afterEach, beforeEach, describe, expect, it } from "bun:test";
import { cleanup, render, screen } from "@solidjs/testing-library";
import { VerifyEmailBlock } from "./VerifyEmailBlock";

beforeEach(() => cleanup());
afterEach(() => cleanup());

describe("VerifyEmailBlock", () => {
  it("shows slice 0 verify copy and resend", () => {
    render(() => <VerifyEmailBlock onResend={() => undefined} resent={false} />);
    expect(screen.getByText("Verify your email to continue.")).toBeTruthy();
    expect(screen.getByRole("button", { name: "Resend the email" })).toBeTruthy();
  });
});
