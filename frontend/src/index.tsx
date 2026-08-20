/** @jsxImportSource solid-js */
import { render } from "solid-js/web";
import App from "./App";
import "./styles/tailwind.css";
import faviconUrl from "./assets/favicon.svg";

const root = document.getElementById("root");

if (!root) {
  throw new Error("Root element not found");
}

root.innerHTML = "";

const link =
  (document.querySelector('link[rel="icon"]') as HTMLLinkElement | null) ??
  document.createElement("link");
link.rel = "icon";
link.href = faviconUrl;
if (!link.parentNode) {
  document.head.appendChild(link);
}

const dispose = render(() => <App />, root);

if (import.meta.hot) {
  import.meta.hot.accept();
  import.meta.hot.dispose(() => {
    dispose();
  });
}
