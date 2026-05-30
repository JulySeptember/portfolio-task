export interface TokenData {
  accessToken: string;
  idToken: string;
  expiresAt: number; // UNIX秒
}

const STORAGE_KEY = "auth_tokens";

export function setTokens(tokens: TokenData) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(tokens));

  // 既存コードとの互換用
  localStorage.setItem("access_token", tokens.accessToken);
  localStorage.setItem("id_token", tokens.idToken);
}

export function getTokens(): TokenData | null {
  const json = localStorage.getItem(STORAGE_KEY);

  if (!json) {
    return null;
  }

  try {
    return JSON.parse(json) as TokenData;
  } catch {
    return null;
  }
}

export function clearTokens() {
  localStorage.removeItem(STORAGE_KEY);

  localStorage.removeItem("access_token");
  localStorage.removeItem("id_token");
}

export function isAccessTokenExpired(): boolean {
  const tokens = getTokens();

  if (!tokens) {
    return true;
  }

  const now = Math.floor(Date.now() / 1000);

  return tokens.expiresAt <= now + 30;
}
