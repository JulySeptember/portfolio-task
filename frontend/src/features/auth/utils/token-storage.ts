const ACCESS_TOKEN_KEY = "access_token";

const ID_TOKEN_KEY = "id_token";

const REFRESH_TOKEN_KEY = "refresh_token";

function isBrowser() {
  return typeof window !== "undefined";
}

export function saveTokens(tokens: {
  access_token: string;
  id_token?: string;
  refresh_token?: string;
}) {
  if (!isBrowser()) {
    return;
  }

  localStorage.setItem(ACCESS_TOKEN_KEY, tokens.access_token);

  if (tokens.id_token) {
    localStorage.setItem(ID_TOKEN_KEY, tokens.id_token);
  }

  if (tokens.refresh_token) {
    localStorage.setItem(REFRESH_TOKEN_KEY, tokens.refresh_token);
  }
}

export function getAccessToken() {
  if (!isBrowser()) {
    return null;
  }

  return localStorage.getItem(ACCESS_TOKEN_KEY);
}

export function getIdToken() {
  if (!isBrowser()) {
    return null;
  }

  return localStorage.getItem(ID_TOKEN_KEY);
}

export function getRefreshToken() {
  if (!isBrowser()) {
    return null;
  }

  return localStorage.getItem(REFRESH_TOKEN_KEY);
}

export function clearTokens() {
  if (!isBrowser()) {
    return;
  }

  localStorage.removeItem(ACCESS_TOKEN_KEY);

  localStorage.removeItem(ID_TOKEN_KEY);

  localStorage.removeItem(REFRESH_TOKEN_KEY);
}
