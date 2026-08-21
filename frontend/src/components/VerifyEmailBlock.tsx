/** @jsxImportSource solid-js */
import { primaryButtonClass } from "./AuthCard";

export function VerifyEmailBlock(props: { onResend: () => void | Promise<void>; resent: boolean }) {
  return (
    <div class="flex flex-col gap-4">
      <p class="text-sm text-ink-700" data-testid="verify-block">
        Verify your email to continue.
      </p>
      <button class={primaryButtonClass} type="button" onClick={props.onResend}>
        Resend the email
      </button>
      {props.resent && <p class="text-sm text-ink-700">Check your email for a link.</p>}
    </div>
  );
}
