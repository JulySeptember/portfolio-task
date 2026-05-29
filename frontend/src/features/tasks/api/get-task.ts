// src/features/tasks/api/get-task.ts

import { apiClient } from "@/lib/api/client";

import { taskSchema, type Task } from "../schemas/task-schema";

import { taskEndpoints } from "./endpoints";

export async function getTask(id: number): Promise<Task> {
  const data = await apiClient<unknown>(taskEndpoints.detail(id));

  return taskSchema.parse(data);
}
