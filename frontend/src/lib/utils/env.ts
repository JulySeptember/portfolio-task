export const env = {
  apiUrl: process.env.NEXT_PUBLIC_API_URL || "",

  cognitoDomain: process.env.NEXT_PUBLIC_COGNITO_DOMAIN || "",

  cognitoClientId: process.env.NEXT_PUBLIC_COGNITO_CLIENT_ID || "",

  cognitoRedirectUri: process.env.NEXT_PUBLIC_COGNITO_REDIRECT_URI || "",

  cognitoLogoutUri: process.env.NEXT_PUBLIC_COGNITO_LOGOUT_URI || "",
};
