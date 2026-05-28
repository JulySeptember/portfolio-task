// src/features/tasks/api/create-task.ts

import { apiClient } from "@/lib/api/client";

import {
  taskSchema,
  type Task,
  type TaskRequest,
} from "../schemas/task-schema";

import { taskEndpoints } from "./endpoints";

export async function createTask(body: TaskRequest): Promise<Task> {
  const data = await apiClient<unknown>(taskEndpoints.list(), {
    method: "POST",

    body: JSON.stringify(body),
  });

  return taskSchema.parse(data);
}
