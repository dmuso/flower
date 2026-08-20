/** @jsxImportSource solid-js */
import { Route, Router } from "@solidjs/router";
import { HomePage } from "./pages/Home";

export default function App() {
  return (
    <Router>
      <Route path="/" component={HomePage} />
    </Router>
  );
}
