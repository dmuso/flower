/** @jsxImportSource solid-js */
import { A, useNavigate } from "@solidjs/router";
import { createSignal } from "solid-js";
import { AuthCard, inputClass, labelClass, primaryButtonClass, secondaryButtonClass } from "../components/AuthCard";
import { ApiRequestError } from "../lib/api/core";
import { requestMagicLink, signIn } from "../lib/api/auth";
import { afterAuthPath } from "../lib/after-auth";

export function SignInPage() {
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
      const me = await signIn(email(), password());
      navigate(afterAuthPath(me));
    } catch (err) {
      if (err instanceof ApiRequestError && err.code === "unauthorized") {
        setError("That email and password don’t match.");
      } else if (err instanceof ApiRequestError) {
        setError(err.message);
      } else {
        setError("That email and password don’t match.");
      }
    } finally {
      setLoading(false);
    }
  }

  return (
    <AuthCard title="Sign in">
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
              autocomplete="current-password"
              value={password()}
              onInput={(e) => setPassword(e.currentTarget.value)}
            />
          </label>
        )}
        {error() && (
          <p class="text-sm text-red-700" data-testid="signin-error">
            {error()}
          </p>
        )}
        <button class={primaryButtonClass} type="submit" disabled={loading()}>
          {loading() ? (mode() === "magic" ? "Sending…" : "Signing in…") : mode() === "magic" ? "Email me a link" : "Sign in"}
        </button>
        <button
          class={secondaryButtonClass}
          type="button"
          onClick={() => {
            setMode(mode() === "password" ? "magic" : "password");
            setError("");
          }}
        >
          {mode() === "password" ? "Email me a link instead" : "Use a password instead"}
        </button>
      </form>
      <p class="text-sm text-ink-700">
        New here?{" "}
        <A class="text-bloom" href="/">
          Sign up
        </A>
      </p>
    </AuthCard>
  );
}
