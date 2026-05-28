// src/features/auth/lib/hosted-ui.ts

export function buildLoginURL() {
  const params = new URLSearchParams({
    client_id: process.env.NEXT_PUBLIC_COGNITO_CLIENT_ID!,

    response_type: "code",

    scope: "email openid profile",

    redirect_uri: process.env.NEXT_PUBLIC_COGNITO_REDIRECT_URI!,
  });

  return `https://${process.env.NEXT_PUBLIC_COGNITO_DOMAIN}/login?${params.toString()}`;
}

export function buildLogoutURL() {
  const params = new URLSearchParams({
    client_id: process.env.NEXT_PUBLIC_COGNITO_CLIENT_ID!,

    logout_uri: process.env.NEXT_PUBLIC_COGNITO_LOGOUT_REDIRECT_URI!,
  });

  return `https://${process.env.NEXT_PUBLIC_COGNITO_DOMAIN}/logout?${params.toString()}`;
}
