import { apiClient } from "@/lib/api/client";

import {
  taskSchema,
  taskStatusSchema,
  type Task,
  type TaskStatus,
} from "../schemas/task-schema";

type UpdateTaskStatusInput = {
  id: number;

  status: TaskStatus;
};

export async function updateTaskStatus({
  id,
  status,
}: UpdateTaskStatusInput): Promise<Task> {
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
