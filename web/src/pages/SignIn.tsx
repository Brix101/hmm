import { SignInInput, signInSchema } from "@/types/auth.type";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import { zodResolver } from "@hookform/resolvers/zod";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { AxiosError } from "axios";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { signInUserMutation } from "@/services/auth.service";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import PasswordInput from "@/components/PasswordInput";
import Icons from "@/components/Icons";
import { ResponseError } from "@/types/error.type";
import { useBoundStore } from "@/store";
import { userSchema } from "@/types/user.type";
import { QUERY_CURRENT_USER_KEY } from "@/constant/query.constant";

function SignIn() {
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const { setPathHistory } = useBoundStore();

  const form = useForm<SignInInput>({
    resolver: zodResolver(signInSchema),
  });

  const { mutate, isLoading } = useMutation({
    mutationFn: signInUserMutation,
    onSuccess: (response) => {
      const user = userSchema.parse(response.data);
      setPathHistory("");
      queryClient.setQueriesData([QUERY_CURRENT_USER_KEY], user);
      navigate("/", { replace: true });
    },
    onError: ({ response }: AxiosError) => {
      const { data } = response as ResponseError;
      for (const fieldName in data.errors) {
        const errorData = data.errors[fieldName];
        form.setError(fieldName as keyof SignInInput, errorData, {
          shouldFocus: true,
        });
      }
    },
  });

  function onSubmit(values: SignInInput) {
    mutate(values);
  }

  return (
    <div className="container flex justify-center items-center h-screen">
      <Card>
        <CardHeader className="space-y-1">
          <CardTitle className="text-2xl">Sign in</CardTitle>
          <CardDescription>
            Choose your preferred sign in method
          </CardDescription>
        </CardHeader>
        <CardContent className="grid gap-4">
          <div className="relative">
            <div className="flex absolute inset-0 items-center">
              <span className="w-full border-t" />
            </div>
            <div className="flex relative justify-center text-xs uppercase">
              {/* <span className="px-2 bg-background text-muted-foreground"> */}
              {/*   Or continue with */}
              {/* </span> */}
            </div>
          </div>
          <div className={cn("grid gap-6 w-96")}>
            <Form {...form}>
              <form
                className="grid gap-4"
                onSubmit={(...args) =>
                  void form.handleSubmit(onSubmit)(...args)
                }
              >
                <FormField
                  control={form.control}
                  name="email"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Email</FormLabel>
                      <FormControl>
                        <Input placeholder="john.doe@example.com" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="password"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Password</FormLabel>
                      <FormControl>
                        <PasswordInput placeholder="**********" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <Button disabled={isLoading}>
                  {isLoading && (
                    <Icons.spinner
                      className="mr-2 w-4 h-4 animate-spin"
                      aria-hidden="true"
                    />
                  )}
                  Sign in
                  <span className="sr-only">Sign in</span>
                </Button>
              </form>
            </Form>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

export default SignIn;
