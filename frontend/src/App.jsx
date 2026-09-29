import React from "react";
import { Nav } from "./components/Nav";
import { ErrorBoundary } from "./components/ErrorBoundary";
import { clearSession, isAuthenticated, isAdmin } from "./state";
import { Feed } from "./pages/Feed";
import { ReportDetail } from "./pages/ReportDetail";
import { NewReport } from "./pages/NewReport";
import { Login } from "./pages/Login";
import { Register } from "./pages/Register";
import { ApplyAuthority } from "./pages/ApplyAuthority";
import { AdminDashboard } from "./pages/admin/AdminDashboard";
import { AdminUsers } from "./pages/admin/AdminUsers";
import { AdminAuthorityRequests } from "./pages/admin/AdminAuthorityRequests";
import { AdminModeration } from "./pages/admin/AdminModeration";

function parseRoute() {
  const raw = window.location.hash.replace(/^#\/?/, "") || "feed";
  const [name, param] = raw.split("/");
  return { name, param };
}

function Forbidden() {
  return <p className="alert error">You do not have access to this page.</p>;
}

export default function App() {
  const [route, setRoute] = React.useState(parseRoute);

  React.useEffect(() => {
    const onChange = () => setRoute(parseRoute());
    window.addEventListener("hashchange", onChange);
    return () => window.removeEventListener("hashchange", onChange);
  }, []);

  const navigate = (path) => {
    window.location.hash = `#/${path}`;
  };
  const logout = () => {
    clearSession();
    navigate("feed");
  };

  let page;
  switch (route.name) {
    case "reports":
      page = <ReportDetail id={route.param} navigate={navigate} />;
      break;
    case "new":
      page = isAuthenticated() ? <NewReport navigate={navigate} /> : <Login navigate={navigate} />;
      break;
    case "login":
      page = <Login navigate={navigate} />;
      break;
    case "register":
      page = <Register navigate={navigate} />;
      break;
    case "apply-authority":
      page = isAuthenticated() ? <ApplyAuthority /> : <Login navigate={navigate} />;
      break;
    case "admin":
      if (!isAdmin()) page = <Forbidden />;
      else if (route.param === "users") page = <AdminUsers />;
      else if (route.param === "authority-requests") page = <AdminAuthorityRequests />;
      else if (route.param === "moderation") page = <AdminModeration />;
      else page = <AdminDashboard />;
      break;
    default:
      page = <Feed />;
  }

  return (
    <>
      <Nav onLogout={logout} />
      <main>
        <ErrorBoundary key={`${route.name}/${route.param || ""}`}>{page}</ErrorBoundary>
      </main>
    </>
  );
}
