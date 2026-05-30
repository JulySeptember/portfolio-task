// src/features/tasks/api/update-task-status.ts

import { apiClient } from "@/lib/api/client";

import { taskSchema, type Task } from "../schemas/task-schema";

import { taskEndpoints } from "./endpoints";

export type UpdateTaskStatusInput = {
  status: "TODO" | "DOING" | "DONE";
};

export async function updateTaskStatus(
  publicId: string,
  input: UpdateTaskStatusInput,
): Promise<Task> {
  const data = await apiClient<unknown>(taskEndpoints.status(publicId), {
    method: "PATCH",

    body: JSON.stringify(input),
  });

  return taskSchema.parse(data);
}
