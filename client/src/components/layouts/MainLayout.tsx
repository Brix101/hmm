import { useGetCurrentUser } from "@/services/user.service";
import { QueryClient } from "@tanstack/react-query";
import { Outlet, useLoaderData } from "react-router-dom";

export const loader =
  ({ queryClient }: { queryClient: QueryClient }) =>
  async () => {
    const query = useGetCurrentUser();

    return (
      queryClient.getQueryData(query.queryKey) ??
      (await queryClient.fetchQuery(query))
    );
  };

function MainLayout() {
  const user = useLoaderData();

  console.log({ user });
  return <Outlet />;
}

export default MainLayout;
