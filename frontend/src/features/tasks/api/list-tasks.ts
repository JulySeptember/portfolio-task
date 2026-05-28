import { apiClient } from "@/lib/api/client";

import {
  taskListResponseSchema,
  type TaskListResponse,
} from "../schemas/task-schema";

export type ListTasksParams = {
  limit?: number;

  offset?: number;

  status?: "TODO" | "DOING" | "DONE";

  sort?: "created_at" | "due_date";

  order?: "ASC" | "DESC";
};

export async function listTasks(
  params?: ListTasksParams,
): Promise<TaskListResponse> {
  const searchParams = new URLSearchParams();

  if (params?.limit) {
    searchParams.set("limit", String(params.limit));
  }

  if (params?.offset) {
    searchParams.set("offset", String(params.offset));
  }

  if (params?.status) {
    searchParams.set("status", params.status);
  }

  if (params?.sort) {
    searchParams.set("sort", params.sort);
  }

  if (params?.order) {
    searchParams.set("order", params.order);
  }

  const query = searchParams.toString();

  const data = await apiClient<unknown>(
    `/api/tasks${query ? `?${query}` : ""}`,
  );

  return taskListResponseSchema.parse(data);
}
