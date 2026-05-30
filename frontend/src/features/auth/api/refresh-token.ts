// src/features/auth/api/refresh-token.ts

import {
  clearTokens,
  getTokens,
  setTokens,
  type TokenData,
} from "@/features/auth/lib/token-storage";

const COGNITO_DOMAIN = process.env.NEXT_PUBLIC_COGNITO_DOMAIN!;
const CLIENT_ID = process.env.NEXT_PUBLIC_COGNITO_CLIENT_ID!;

let refreshPromise: Promise<TokenData | null> | null = null;

async function doRefresh(): Promise<TokenData | null> {
  const tokens = getTokens();

  if (!tokens?.refreshToken) {
    return null;
  }

  const body = new URLSearchParams({
    grant_type: "refresh_token",
    client_id: CLIENT_ID,
    refresh_token: tokens.refreshToken,
  });

  const response = await fetch(`https://${COGNITO_DOMAIN}/oauth2/token`, {
    method: "POST",
    headers: {
      "Content-Type": "application/x-www-form-urlencoded",
    },
    body: body.toString(),
  });

  if (!response.ok) {
    clearTokens();
    return null;
  }

  const data = await response.json();

  const refreshedTokens: TokenData = {
    accessToken: data.access_token,
    idToken: data.id_token ?? tokens.idToken,
    refreshToken: tokens.refreshToken,
    expiresAt: Math.floor(Date.now() / 1000) + Number(data.expires_in),
  };

  setTokens(refreshedTokens);

  return refreshedTokens;
}

export async function refreshAccessToken(): Promise<TokenData | null> {
  if (refreshPromise) {
    return refreshPromise;
  }

  refreshPromise = doRefresh();

  try {
    return await refreshPromise;
  } finally {
    refreshPromise = null;
  }
}
