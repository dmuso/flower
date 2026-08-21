/** @jsxImportSource solid-js */
import { Route, Router } from "@solidjs/router";
import { SignUpPage } from "./pages/SignUp";
import { SignInPage } from "./pages/SignIn";
import { CheckEmailPage } from "./pages/CheckEmail";
import { VerifyEmailPage } from "./pages/VerifyEmail";
import { MagicLinkPage } from "./pages/MagicLink";
import { NameOrganisationPage } from "./pages/NameOrganisation";
import { NameProjectPage } from "./pages/NameProject";
import { BoardPage } from "./pages/Board";

export default function App() {
  return (
    <Router>
      <Route path="/" component={SignUpPage} />
      <Route path="/signin" component={SignInPage} />
      <Route path="/check-email" component={CheckEmailPage} />
      <Route path="/verify-email" component={VerifyEmailPage} />
      <Route path="/magic-link" component={MagicLinkPage} />
      <Route path="/organisations/new" component={NameOrganisationPage} />
      <Route path="/organisations/:orgId/projects/new" component={NameProjectPage} />
      <Route path="/projects/:projectId" component={BoardPage} />
    </Router>
  );
}
