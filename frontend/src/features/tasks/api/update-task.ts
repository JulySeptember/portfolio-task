// src/features/tasks/api/update-task.ts

import { apiClient } from "@/lib/api/client";

import {
  taskSchema,
  type Task,
  type TaskRequest,
} from "../schemas/task-schema";

import { taskEndpoints } from "./endpoints";

export async function updateTask(id: number, body: TaskRequest): Promise<Task> {
  const data = await apiClient<unknown>(taskEndpoints.detail(id), {
    method: "PUT",

    body: JSON.stringify(body),
  });

  return taskSchema.parse(data);
}
