import { getRefreshToken, saveTokens } from "../utils/token-storage";

import { refreshTokenResponseSchema } from "../types/auth";

export async function refreshToken() {
  const refreshToken = getRefreshToken();

  if (!refreshToken) {
    throw new Error("No refresh token");
  }

  const body = new URLSearchParams({
    grant_type: "refresh_token",

    client_id: process.env.NEXT_PUBLIC_COGNITO_CLIENT_ID!,

    refresh_token: refreshToken,
  });

  const response = await fetch(
    `https://${process.env.NEXT_PUBLIC_COGNITO_DOMAIN}/oauth2/token`,
    {
      method: "POST",

      headers: {
        "Content-Type": "application/x-www-form-urlencoded",
      },

      body,
    },
  );

  if (!response.ok) {
    throw new Error("Refresh failed");
  }

  const json = await response.json();

  const tokens = refreshTokenResponseSchema.parse(json);

  saveTokens(tokens);

  return tokens.access_token;
}
