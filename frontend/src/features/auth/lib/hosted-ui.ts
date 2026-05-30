// src/features/auth/lib/hosted-ui.ts

const COGNITO_DOMAIN = process.env.NEXT_PUBLIC_COGNITO_DOMAIN!.replace(
  /^https?:\/\//,
  "",
);

const CLIENT_ID = process.env.NEXT_PUBLIC_COGNITO_CLIENT_ID!;
const REDIRECT_URI = process.env.NEXT_PUBLIC_COGNITO_REDIRECT_URI!;
const LOGOUT_REDIRECT_URI =
  process.env.NEXT_PUBLIC_COGNITO_LOGOUT_REDIRECT_URI!;

export async function buildLoginURL(): Promise<string> {
  const params = new URLSearchParams({
    client_id: CLIENT_ID,
    response_type: "token",
    scope: "openid email profile",
    redirect_uri: REDIRECT_URI,
  });

  return `https://${COGNITO_DOMAIN}/login?${params.toString()}`;
}

export function buildLogoutURL(): string {
  const params = new URLSearchParams({
    client_id: CLIENT_ID,
    logout_uri: LOGOUT_REDIRECT_URI,
  });

  return `https://${COGNITO_DOMAIN}/logout?${params.toString()}`;
}
