// src/features/tasks/api/update-task.ts

import { apiClient } from "@/lib/api/client";

import { taskSchema, type Task } from "../schemas/task-schema";

import { taskEndpoints } from "./endpoints";

export type UpdateTaskInput = {
  title: string;

  description?: string;

  status: "TODO" | "DOING" | "DONE";

  due_date?: string | null;
};

export async function updateTask(
  publicId: string,
  input: UpdateTaskInput,
): Promise<Task> {
  const data = await apiClient<unknown>(taskEndpoints.detail(publicId), {
    method: "PUT",

    body: JSON.stringify(input),
  });

  return taskSchema.parse(data);
}
