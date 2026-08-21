/** @jsxImportSource solid-js */
import { useSearchParams } from "@solidjs/router";
import { AuthCard } from "../components/AuthCard";

export function CheckEmailPage() {
  const [params] = useSearchParams();
  const magic = () => params.kind === "magic";
  return (
    <AuthCard title="Check your email">
      <p class="text-sm text-ink-700" data-testid="check-email-copy">
        {magic() ? "Check your email for a link." : "Check your email to verify."}
      </p>
    </AuthCard>
  );
}
