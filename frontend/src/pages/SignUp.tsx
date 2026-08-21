/** @jsxImportSource solid-js */
import { A, useNavigate } from "@solidjs/router";
import { createSignal } from "solid-js";
import { AuthCard, inputClass, labelClass, primaryButtonClass, secondaryButtonClass } from "../components/AuthCard";
import { ApiRequestError } from "../lib/api/core";
import { requestMagicLink, signUp } from "../lib/api/auth";

export function SignUpPage() {
  const navigate = useNavigate();
  const [email, setEmail] = createSignal("");
  const [password, setPassword] = createSignal("");
  const [mode, setMode] = createSignal<"password" | "magic">("password");
  const [error, setError] = createSignal("");
  const [loading, setLoading] = createSignal(false);

  async function onSubmit(event: Event) {
    event.preventDefault();
    setError("");
    setLoading(true);
    try {
      if (mode() === "magic") {
        await requestMagicLink(email());
        navigate("/check-email?kind=magic");
        return;
      }
      await signUp(email(), password());
      navigate("/check-email");
    } catch (err) {
      if (err instanceof ApiRequestError && err.status === 409) {
        setError("That email already belongs to a user.");
      } else if (err instanceof ApiRequestError) {
        setError(err.message);
      } else {
        setError("Need an email and a password.");
      }
    } finally {
      setLoading(false);
    }
  }

  return (
    <AuthCard title="Sign up">
      <p class="text-sm text-ink-700">Email and a password. That’s enough to start.</p>
      <form class="flex flex-col gap-4" onSubmit={onSubmit}>
        <label class="flex flex-col gap-1">
          <span class={labelClass}>Email</span>
          <input
            class={inputClass}
            type="email"
            name="email"
            autocomplete="email"
            value={email()}
            onInput={(e) => setEmail(e.currentTarget.value)}
          />
        </label>
        {mode() === "password" && (
          <label class="flex flex-col gap-1">
            <span class={labelClass}>Password</span>
            <input
              class={inputClass}
              type="password"
              name="password"
              autocomplete="new-password"
              value={password()}
              onInput={(e) => setPassword(e.currentTarget.value)}
            />
          </label>
        )}
        {error() && (
          <p class="text-sm text-red-700" data-testid="signup-error">
            {error()}
          </p>
        )}
        <button class={primaryButtonClass} type="submit" disabled={loading()}>
          {loading() ? (mode() === "magic" ? "Sending…" : "Signing you up…") : mode() === "magic" ? "Email me a link" : "Sign up"}
        </button>
        <button
          class={secondaryButtonClass}
          type="button"
          onClick={() => {
            setMode(mode() === "password" ? "magic" : "password");
            setError("");
          }}
        >
          Email me a link instead
        </button>
      </form>
      <p class="text-sm text-ink-700">
        Already a user?{" "}
        <A class="text-bloom" href="/signin">
          Sign in
        </A>
      </p>
    </AuthCard>
  );
}
