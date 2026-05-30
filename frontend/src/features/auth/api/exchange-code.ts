import { setTokens } from "@/features/auth/lib/token-storage";

const COGNITO_DOMAIN = process.env.NEXT_PUBLIC_COGNITO_DOMAIN!;
const CLIENT_ID = process.env.NEXT_PUBLIC_COGNITO_CLIENT_ID!;
const REDIRECT_URI = process.env.NEXT_PUBLIC_COGNITO_REDIRECT_URI!;

export async function exchangeCodeForTokens(code: string) {
  const codeVerifier = sessionStorage.getItem("pkce_code_verifier");
  if (!codeVerifier) throw new Error("No PKCE code_verifier found");

  const body = new URLSearchParams({
    grant_type: "authorization_code",
    client_id: CLIENT_ID,
    code,
    redirect_uri: REDIRECT_URI,
    code_verifier: codeVerifier,
  });

  const res = await fetch(`https://${COGNITO_DOMAIN}/oauth2/token`, {
    method: "POST",
    headers: { "Content-Type": "application/x-www-form-urlencoded" },
    body: body.toString(),
  });

  if (!res.ok) throw new Error(`Token exchange failed: ${await res.text()}`);

  const data = await res.json();

  const expiresAt = Math.floor(Date.now() / 1000) + data.expires_in;

  setTokens({
    accessToken: data.access_token,
    idToken: data.id_token,
    refreshToken: data.refresh_token,
    expiresAt,
  });

  sessionStorage.removeItem("pkce_code_verifier");
}
