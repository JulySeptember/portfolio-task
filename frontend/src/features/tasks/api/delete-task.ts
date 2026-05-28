// src/features/tasks/api/delete-task.ts

import { apiClient } from "@/lib/api/client";

import { taskEndpoints } from "./endpoints";

export async function deleteTask(id: number): Promise<void> {
  await apiClient(taskEndpoints.detail(id), {
    method: "DELETE",
  });
}
