import { apiClient } from "@/lib/api/client";

import {
  taskSchema,
  taskRequestSchema,
  type Task,
  type TaskRequest,
} from "../schemas/task-schema";

export async function updateTask(
  id: number,
  input: TaskRequest,
): Promise<Task> {
  const body = taskRequestSchema.parse(input);

  const data = await apiClient<unknown>(
    `${process.env.NEXT_PUBLIC_API_URL}/tasks/${id}`,
    {
      method: "PUT",

      body: JSON.stringify(body),
    },
  );

  return taskSchema.parse(data);
}
