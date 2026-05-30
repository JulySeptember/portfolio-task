import { apiClient } from "@/lib/api/client";
import { taskEndpoints } from "./endpoints";

export async function deleteTask(publicId: string): Promise<void> {
  await apiClient(taskEndpoints.delete(publicId), { method: "DELETE" });
}
