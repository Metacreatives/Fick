import { Link, Outlet } from "react-router";

export function RootLayout() {
  return (
    <>
      <header>
        <Link to="/">Fick</Link>
      </header>

      <main>
        <Outlet />
      </main>
    </>
  );
}
