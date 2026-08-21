/** @jsxImportSource solid-js */
import { useNavigate, useParams } from "@solidjs/router";
import { createSignal } from "solid-js";
import { AuthCard, inputClass, labelClass, primaryButtonClass } from "../components/AuthCard";
import { ApiRequestError } from "../lib/api/core";
import { createProject } from "../lib/api/projects";

export function NameProjectPage() {
  const navigate = useNavigate();
  const params = useParams();
  const [name, setName] = createSignal("");
  const [error, setError] = createSignal("");
  const [loading, setLoading] = createSignal(false);

  async function onSubmit(event: Event) {
    event.preventDefault();
    setError("");
    setLoading(true);
    try {
      const project = await createProject(params.orgId, name());
      navigate(`/projects/${project.id}`);
    } catch (err) {
      if (err instanceof ApiRequestError) {
        setError(err.message);
      } else {
        setError("Name the project.");
      }
    } finally {
      setLoading(false);
    }
  }

  return (
    <AuthCard title="Name the first project">
      <p class="text-sm text-ink-700">Name the first project.</p>
      <form class="flex flex-col gap-4" onSubmit={onSubmit}>
        <label class="flex flex-col gap-1">
          <span class={labelClass}>Project name</span>
          <input
            class={inputClass}
            type="text"
            name="project"
            value={name()}
            onInput={(e) => setName(e.currentTarget.value)}
          />
        </label>
        {error() && (
          <p class="text-sm text-red-700" data-testid="project-error">
            {error()}
          </p>
        )}
        <button class={primaryButtonClass} type="submit" disabled={loading()}>
          Create project
        </button>
      </form>
    </AuthCard>
  );
}
