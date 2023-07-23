import { QUERY_CURRENT_USER_KEY } from "@/constant/query.constant";
import { apiClient } from "@/lib/httpCommon";
import { userSchema } from "@/types/user.type";

export const useGetCurrentUser = () => ({
  queryKey: [QUERY_CURRENT_USER_KEY],
  queryFn: async () => {
    const res = await apiClient.get("/me");
    return userSchema.parse(res.data);
  },
});
