import { apiClient } from "@/lib/api/client";

import {
  taskSchema,
  type Task,
  type TaskRequest,
} from "../schemas/task-schema";

type Params = {
  id: number;
} & TaskRequest;

export async function updateTask({ id, ...body }: Params): Promise<Task> {
  const data = await apiClient<unknown>(`/api/tasks/${id}`, {
    method: "PUT",

    body: JSON.stringify(body),
  });

  return taskSchema.parse(data);
}
