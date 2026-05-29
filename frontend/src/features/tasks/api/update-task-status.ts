// src/features/tasks/api/update-task-status.ts

import { apiClient } from "@/lib/api/client";

import { taskSchema, type Task, type TaskStatus } from "../schemas/task-schema";

import { taskEndpoints } from "./endpoints";

type UpdateTaskStatusBody = {
  status: TaskStatus;
};

export async function updateTaskStatus(
  id: number,
  body: UpdateTaskStatusBody,
): Promise<Task> {
  const data = await apiClient<unknown>(taskEndpoints.status(id), {
    method: "PATCH",

    body: JSON.stringify(body),
  });

  return taskSchema.parse(data);
}
