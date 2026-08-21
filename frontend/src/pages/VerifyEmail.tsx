/** @jsxImportSource solid-js */
import { useNavigate, useSearchParams } from "@solidjs/router";
import { createSignal, onMount } from "solid-js";
import { AuthCard, primaryButtonClass } from "../components/AuthCard";
import { consumeVerifyEmail } from "../lib/api/auth";
import { afterAuthPath } from "../lib/after-auth";

export function VerifyEmailPage() {
  const navigate = useNavigate();
  const [params] = useSearchParams();
  const [error, setError] = createSignal("");

  onMount(async () => {
    const token = params.token;
    if (!token || Array.isArray(token)) {
      setError("That link is no longer valid.");
      return;
    }
    try {
      const me = await consumeVerifyEmail(token);
      navigate(afterAuthPath(me));
    } catch {
      setError("That link is no longer valid.");
    }
  });

  return (
    <AuthCard title="Verify your email">
      {error() ? (
        <div class="flex flex-col gap-4">
          <p class="text-sm text-red-700">{error()}</p>
          <a class={primaryButtonClass + " inline-flex w-fit"} href="/">
            Sign up
          </a>
        </div>
      ) : (
        <p class="text-sm text-ink-700">Verify your email to continue.</p>
      )}
    </AuthCard>
  );
}
