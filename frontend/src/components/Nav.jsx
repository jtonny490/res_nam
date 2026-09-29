import React from "react";
import { isAuthenticated, isAdmin, isPublic } from "../state";

export function Nav({ onLogout }) {
  const link = (path, label) => <a href={`#/${path}`}>{label}</a>;
  return (
    <nav className="app-nav">
      <span className="brand">RES NAM</span>
      {link("feed", "Reports")}
      {isAuthenticated() && link("new", "New report")}
      {isAdmin() && link("admin", "Admin")}
      {isAuthenticated() && isPublic() && link("apply-authority", "Apply for authority")}
      <span className="spacer" />
      {isAuthenticated() ? (
        <button type="button" onClick={onLogout}>
          Log out
        </button>
      ) : (
        <>
          {link("login", "Log in")}
          {link("register", "Register")}
        </>
      )}
    </nav>
  );
}
