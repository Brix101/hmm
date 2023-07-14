import React, { Suspense } from "react";
import { createBrowserRouter } from "react-router-dom";

const About = React.lazy(() => import("@/pages/About"));
const Home = React.lazy(() => import("@/pages/Home"));
const Contact = React.lazy(() => import("@/pages/Contact"));

const router = createBrowserRouter([
  {
    path: "/",
    element: (
      <Suspense fallback={"Loading..."}>
        <Home />
      </Suspense>
    ),
  },
  {
    path: "/about",
    element: (
      <Suspense fallback={"Loading..."}>
        <About />
      </Suspense>
    ),
  },
  {
    path: "/contact",
    element: (
      <Suspense fallback={"Loading..."}>
        <Contact />
      </Suspense>
    ),
  },
]);

export default router;
