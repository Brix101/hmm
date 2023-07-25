import { apiClient } from "@/lib/httpCommon";
import { SignInInput } from "@/types/auth.type";

function signInUserMutation({ email, password }: SignInInput) {
  return apiClient.post(
    "users/sign-in",
    JSON.stringify({ email: email, password: password }),
    {
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
    }
  );
}

export { signInUserMutation };
