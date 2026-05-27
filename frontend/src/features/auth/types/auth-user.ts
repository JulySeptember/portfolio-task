import { z } from "zod";

export const currentUserSchema = z.object({
  id: z.number(),

  email: z.email(),

  name: z.string(),
});

export type CurrentUser = z.infer<typeof currentUserSchema>;
