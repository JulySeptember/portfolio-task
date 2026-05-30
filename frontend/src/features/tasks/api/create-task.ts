// src/features/tasks/api/create-task.ts

import { apiClient } from "@/lib/api/client";

import { taskSchema, type Task } from "../schemas/task-schema";

import { taskEndpoints } from "./endpoints";

export type CreateTaskInput = {
  title: string;

  description?: string;

  status: "TODO" | "DOING" | "DONE";

  due_date?: string | null;
};

export async function createTask(input: CreateTaskInput): Promise<Task> {
  const data = await apiClient<unknown>(taskEndpoints.list(), {
    method: "POST",

    body: JSON.stringify(input),
  });

  return taskSchema.parse(data);
}
