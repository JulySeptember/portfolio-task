// src/features/auth/lib/token.ts

export function getAccessToken() {
  if (typeof window === "undefined") {
    return null;
  }

  return localStorage.getItem("access_token");
}

export function getIdToken() {
  if (typeof window === "undefined") {
    return null;
  }

  return localStorage.getItem("id_token");
}

export function getRefreshToken() {
  if (typeof window === "undefined") {
    return null;
  }

  return localStorage.getItem("refresh_token");
}

export function setTokens(tokens: {
  accessToken: string;
  idToken: string;
  refreshToken?: string;
}) {
  localStorage.setItem("access_token", tokens.accessToken);

  localStorage.setItem("id_token", tokens.idToken);

  if (tokens.refreshToken) {
    localStorage.setItem("refresh_token", tokens.refreshToken);
  }
}

export function clearTokens() {
  localStorage.removeItem("access_token");

  localStorage.removeItem("id_token");

  localStorage.removeItem("refresh_token");
}

export function isAuthenticated() {
  return !!getAccessToken();
}
