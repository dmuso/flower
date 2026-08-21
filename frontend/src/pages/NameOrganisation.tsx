/** @jsxImportSource solid-js */
import { useNavigate } from "@solidjs/router";
import { createSignal, onMount } from "solid-js";
import { AuthCard, inputClass, labelClass, primaryButtonClass } from "../components/AuthCard";
import { VerifyEmailBlock } from "../components/VerifyEmailBlock";
import { ApiRequestError } from "../lib/api/core";
import { fetchMe, requestMagicLink } from "../lib/api/auth";
import { createOrganisation } from "../lib/api/organisations";

export function NameOrganisationPage() {
  const navigate = useNavigate();
  const [name, setName] = createSignal("");
  const [error, setError] = createSignal("");
  const [loading, setLoading] = createSignal(false);
  const [blocked, setBlocked] = createSignal("");
  const [email, setEmail] = createSignal("");
  const [resent, setResent] = createSignal(false);

  onMount(async () => {
    try {
      const me = await fetchMe();
      if (!me.email_verified_at) {
        setEmail(me.email);
        setBlocked("Verify your email to continue.");
      }
    } catch (err) {
      if (err instanceof ApiRequestError && err.status === 401) {
        navigate("/signin");
      }
    }
  });

  async function onSubmit(event: Event) {
    event.preventDefault();
    setError("");
    setLoading(true);
    try {
      const org = await createOrganisation(name());
      navigate(`/organisations/${org.id}/projects/new`);
    } catch (err) {
      if (err instanceof ApiRequestError && err.code === "email_unverified") {
        setBlocked("Verify your email to continue.");
      } else if (err instanceof ApiRequestError) {
        setError(err.message);
      } else {
        setError("Name the organisation.");
      }
    } finally {
      setLoading(false);
    }
  }

  return (
    <AuthCard title="Name the organisation">
      {blocked() ? (
        <VerifyEmailBlock
          resent={resent()}
          onResend={async () => {
            await requestMagicLink(email());
            setResent(true);
          }}
        />
      ) : (
        <>
          <p class="text-sm text-ink-700">What’s the organisation called?</p>
          <form class="flex flex-col gap-4" onSubmit={onSubmit}>
            <label class="flex flex-col gap-1">
              <span class={labelClass}>Organisation name</span>
              <input
                class={inputClass}
                type="text"
                name="organisation"
                value={name()}
                onInput={(e) => setName(e.currentTarget.value)}
              />
            </label>
            {error() && (
              <p class="text-sm text-red-700" data-testid="org-error">
                {error()}
              </p>
            )}
            <button class={primaryButtonClass} type="submit" disabled={loading()}>
              Create organisation
            </button>
          </form>
        </>
      )}
    </AuthCard>
  );
}
