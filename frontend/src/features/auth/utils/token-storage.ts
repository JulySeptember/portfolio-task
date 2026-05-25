const ACCESS_TOKEN_KEY =
  "access_token"

const ID_TOKEN_KEY =
  "id_token"

const REFRESH_TOKEN_KEY =
  "refresh_token"

export function saveTokens(tokens: {
  access_token: string
  id_token: string
  refresh_token?: string
}) {
  localStorage.setItem(
    ACCESS_TOKEN_KEY,
    tokens.access_token
  )

  localStorage.setItem(
    ID_TOKEN_KEY,
    tokens.id_token
  )

  if (tokens.refresh_token) {
    localStorage.setItem(
      REFRESH_TOKEN_KEY,
      tokens.refresh_token
    )
  }
}

export function getAccessToken() {
  return localStorage.getItem(
    ACCESS_TOKEN_KEY
  )
}

export function getIdToken() {
  return localStorage.getItem(
    ID_TOKEN_KEY
  )
}

export function clearTokens() {
  localStorage.removeItem(
    ACCESS_TOKEN_KEY
  )

  localStorage.removeItem(
    ID_TOKEN_KEY
  )

  localStorage.removeItem(
    REFRESH_TOKEN_KEY
  )
}
