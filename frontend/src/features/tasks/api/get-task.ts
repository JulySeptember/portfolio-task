import { apiClient } from "@/lib/api/client";

import { taskSchema, type Task } from "../schemas/task-schema";

export async function getTask(id: number): Promise<Task> {
  const data = await apiClient<unknown>(`/api/tasks/${id}`);

  return taskSchema.parse(data);
}
