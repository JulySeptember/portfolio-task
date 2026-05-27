import { apiClient } from "@/lib/api/client";

import { taskSchema, type Task } from "../schemas/task-schema";

type Params = {
  id: number;

  status: "TODO" | "DOING" | "DONE";
};

export async function updateTaskStatus({ id, status }: Params): Promise<Task> {
  const data = await apiClient<unknown>(`/api/tasks/${id}/status`, {
    method: "PATCH",

    body: JSON.stringify({
      status,
    }),
  });

  return taskSchema.parse(data);
}
