import MainLayout, {
  loader as userLoader,
} from "@/components/layouts/MainLayout";
import { QueryClient } from "@tanstack/react-query";
import React, { Suspense } from "react";
import { createBrowserRouter } from "react-router-dom";

const About = React.lazy(() => import("@/pages/About"));
const Home = React.lazy(() => import("@/pages/Home"));
const Contact = React.lazy(() => import("@/pages/Contact"));
const SignIn = React.lazy(() => import("@/pages/SignIn"));

const queryClient = new QueryClient();

const router = createBrowserRouter([
  {
    path: "/",
    element: <MainLayout />,
    loader: userLoader({ queryClient }),
    errorElement: <SignIn />,
    children: [
      {
        index: true,
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
    ],
  },
]);

export default router;
