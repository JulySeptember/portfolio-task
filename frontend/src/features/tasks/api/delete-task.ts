import { apiClient } from "@/lib/api/client";

export async function deleteTask(id: number): Promise<void> {
  await apiClient(`${process.env.NEXT_PUBLIC_API_URL}/tasks/${id}`, {
    method: "DELETE",
  });
}
