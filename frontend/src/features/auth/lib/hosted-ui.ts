// src/features/auth/lib/hosted-ui.ts

export async function buildLoginURL(): Promise<string> {
  const domain = process.env.NEXT_PUBLIC_COGNITO_DOMAIN;
  const clientId = process.env.NEXT_PUBLIC_COGNITO_CLIENT_ID;
  const redirectUri = process.env.NEXT_PUBLIC_COGNITO_REDIRECT_URI;

  if (!domain || !clientId || !redirectUri) {
    throw new Error("Cognito environment variables are not configured");
  }

  const params = new URLSearchParams({
    client_id: clientId,
    response_type: "token",
    scope: "openid email profile",
    redirect_uri: redirectUri,
  });

  return `https://${domain.replace(/^https?:\/\//, "")}/login?${params.toString()}`;
}

export function buildLogoutURL(): string {
  const domain = process.env.NEXT_PUBLIC_COGNITO_DOMAIN;
  const clientId = process.env.NEXT_PUBLIC_COGNITO_CLIENT_ID;
  const logoutRedirectUri = process.env.NEXT_PUBLIC_COGNITO_LOGOUT_REDIRECT_URI;

  if (!domain || !clientId || !logoutRedirectUri) {
    throw new Error("Cognito environment variables are not configured");
  }

  const params = new URLSearchParams({
    client_id: clientId,
    logout_uri: logoutRedirectUri,
  });

  return `https://${domain.replace(/^https?:\/\//, "")}/logout?${params.toString()}`;
}
