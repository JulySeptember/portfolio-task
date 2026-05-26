import { apiClient } from "@/lib/api/client";

import {
  taskSchema,
  type Task,
  type TaskRequest,
} from "../schemas/task-schema";

type UpdateTaskInput = {
  id: number;
} & TaskRequest;

export async function updateTask({
  id,
  title,
  description,
  status,
  due_date,
}: UpdateTaskInput): Promise<Task> {
  const data = await apiClient<unknown>(
    `${process.env.NEXT_PUBLIC_API_URL}/tasks/${id}`,
    {
      method: "PUT",

      body: JSON.stringify({
        title,
        description,
        status,
        due_date,
      }),
    },
  );

  return taskSchema.parse(data);
}
