import { apiClient } from "@/lib/api/client";

export async function deleteTask(id: number): Promise<void> {
  await apiClient(`/api/tasks/${id}`, {
    method: "DELETE",
  });
}
