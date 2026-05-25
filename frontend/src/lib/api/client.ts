export async function apiClient(
  path: string,
  options?: RequestInit
) {
  const baseUrl =
    process.env.NEXT_PUBLIC_API_URL!

  const response =
    await fetch(
      `${baseUrl}${path}`,
      {
        ...options,
        headers: {
          "Content-Type":
            "application/json",

          ...(options?.headers || {}),
        },
      }
    )

  if (!response.ok) {
    throw new Error(
      "api request failed"
    )
  }

  if (response.status === 204) {
    return null
  }

  return response.json()
}
