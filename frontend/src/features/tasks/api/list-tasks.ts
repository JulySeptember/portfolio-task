import { apiClient } from "@/lib/api/client";

import {
  taskListResponseSchema,
  type TaskListResponse,
} from "../schemas/task-schema";

import { taskEndpoints } from "./endpoints";

export type ListTasksParams = {
  limit?: number;
  offset?: number;
  status?: "TODO" | "DOING" | "DONE";
  sort?: "created_at" | "due_date";
  order?: "ASC" | "DESC";
};

export async function listTasks(
  params: ListTasksParams = {},
): Promise<TaskListResponse> {
  const searchParams = new URLSearchParams();

  if (params.limit !== undefined) {
    searchParams.set("limit", String(params.limit));
  }

  if (params.offset !== undefined) {
    searchParams.set("offset", String(params.offset));
  }

  if (params.status) {
    searchParams.set("status", params.status);
  }

  if (params.sort) {
    searchParams.set("sort", params.sort);
  }

  if (params.order) {
    searchParams.set("order", params.order);
  }

  const endpoint = searchParams.toString()
    ? `${taskEndpoints.list()}?${searchParams.toString()}`
    : taskEndpoints.list();

  const data = await apiClient<unknown>(endpoint);

  return taskListResponseSchema.parse(data);
}
