/** @jsxImportSource solid-js */
import { afterEach, beforeEach, describe, expect, it } from "bun:test";
import { cleanup, render, screen } from "@solidjs/testing-library";
import { primaryButtonClass } from "./AuthCard";
import { VerifyEmailBlock } from "./VerifyEmailBlock";

beforeEach(() => cleanup());
afterEach(() => cleanup());

describe("VerifyEmailBlock", () => {
  it("shows verify-your-email copy and resend as the primary next action", () => {
    render(() => <VerifyEmailBlock onResend={() => undefined} resent={false} />);
    expect(screen.getByText("Verify your email to continue.")).toBeTruthy();
    const resend = screen.getByRole("button", { name: "Resend the email" });
    expect(resend).toBeTruthy();
    expect(resend.getAttribute("class")).toBe(primaryButtonClass);
  });
});
