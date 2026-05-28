// features/auth/types/auth.ts

import { z } from "zod";

export const refreshTokenResponseSchema = z.object({
  access_token: z.string(),

  id_token: z.string().optional(),

  refresh_token: z.string().optional(),

  expires_in: z.number(),

  token_type: z.string(),
});

export type RefreshTokenResponse = z.infer<typeof refreshTokenResponseSchema>;
