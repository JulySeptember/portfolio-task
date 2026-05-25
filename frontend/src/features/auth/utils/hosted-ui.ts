const DOMAIN =
  process.env.NEXT_PUBLIC_COGNITO_DOMAIN!

const CLIENT_ID =
  process.env.NEXT_PUBLIC_COGNITO_CLIENT_ID!

const REDIRECT_URI =
  process.env.NEXT_PUBLIC_COGNITO_REDIRECT_URI!

const LOGOUT_URI =
  process.env.NEXT_PUBLIC_COGNITO_LOGOUT_URI!

export function buildLoginURL() {
  const params =
    new URLSearchParams({
      client_id:
        CLIENT_ID,

      response_type:
        "code",

      scope:
        "openid email profile",

      redirect_uri:
        REDIRECT_URI,
    })

  return `https://${DOMAIN}/login?${params.toString()}`
}

export function buildLogoutURL() {
  const params =
    new URLSearchParams({
      client_id:
        CLIENT_ID,

      logout_uri:
        LOGOUT_URI,
    })

  return `https://${DOMAIN}/logout?${params.toString()}`
}
