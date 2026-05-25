export async function fetcher<T>(
  input: RequestInfo,
  init?: RequestInit
): Promise<T> {
  const response =
    await fetch(input, init)

  if (!response.ok) {
    throw new Error("fetch failed")
  }

  return response.json()
}
