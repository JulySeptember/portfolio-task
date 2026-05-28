import { apiClient } from "@/lib/api/client";

import {
  taskSchema,
  type Task,
  type TaskRequest,
} from "../schemas/task-schema";

export async function createTask(body: TaskRequest): Promise<Task> {
  const data = await apiClient<unknown>("/api/tasks", {
    method: "POST",

    body: JSON.stringify(body),
  });

  return taskSchema.parse(data);
}
