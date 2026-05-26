import { apiClient } from "@/lib/api/client";

import { taskSchema, type Task } from "../schemas/task-schema";

export async function getTask(id: number): Promise<Task> {
  const data = await apiClient<unknown>(
    `${process.env.NEXT_PUBLIC_API_URL}/tasks/${id}`,
  );

  return taskSchema.parse(data);
}
