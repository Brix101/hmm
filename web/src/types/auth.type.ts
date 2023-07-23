import z from "zod";

export const signInSchema = z.object({
  email: z.string().email(),
  password: z
    .string()
    .min(8, { message: "Password must containt atleast 8 character(s)" }),
});

export type SignInInput = z.TypeOf<typeof signInSchema>;
