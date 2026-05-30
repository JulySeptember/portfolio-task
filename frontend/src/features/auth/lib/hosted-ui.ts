// src/features/auth/lib/hosted-ui.ts

import { generatePKCE } from "./pkce";

const COGNITO_DOMAIN = process.env.NEXT_PUBLIC_COGNITO_DOMAIN!;
const CLIENT_ID = process.env.NEXT_PUBLIC_COGNITO_CLIENT_ID!;
const REDIRECT_URI = process.env.NEXT_PUBLIC_COGNITO_REDIRECT_URI!;
const LOGOUT_REDIRECT_URI =
  process.env.NEXT_PUBLIC_COGNITO_LOGOUT_REDIRECT_URI!;

/**
 * Build Cognito Login URL (Authorization Code Flow + PKCE)
 */
export async function buildLoginURL(): Promise<string> {
  const { codeChallenge } = await generatePKCE();

  const params = new URLSearchParams({
    client_id: CLIENT_ID,
    response_type: "code",
    scope: "openid email profile",
    redirect_uri: REDIRECT_URI,
    code_challenge: codeChallenge,
    code_challenge_method: "S256",
  });

  return `https://${COGNITO_DOMAIN}/login?${params.toString()}`;
}

/**
 * Build Cognito Logout URL
 */
export function buildLogoutURL(): string {
  const params = new URLSearchParams({
    client_id: CLIENT_ID,
    logout_uri: LOGOUT_REDIRECT_URI,
  });

  return `https://${COGNITO_DOMAIN}/logout?${params.toString()}`;
}
