import { apiClient } from "@/lib/api/client";

import {
  taskSchema,
  taskStatusSchema,
  type Task,
  type TaskStatus,
} from "../schemas/task-schema";

export async function updateTaskStatus(
  id: number,
  status: TaskStatus,
): Promise<Task> {
  const body = {
    status: taskStatusSchema.parse(status),
  };

  const data = await apiClient<unknown>(
    `${process.env.NEXT_PUBLIC_API_URL}/tasks/${id}/status`,
    {
      method: "PATCH",

      body: JSON.stringify(body),
    },
  );

  return taskSchema.parse(data);
}
