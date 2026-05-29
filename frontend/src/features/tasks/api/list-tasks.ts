// src/features/tasks/api/list-tasks.ts

import { apiClient } from "@/lib/api/client";

import {
  taskListResponseSchema,
  type TaskListResponse,
} from "../schemas/task-schema";

import { taskEndpoints } from "./endpoints";

export type ListTasksParams = {
  limit?: number;

  offset?: number;

  status?: string;

  sort?: string;

  order?: string;
};

export async function listTasks(
  params: ListTasksParams = {},
): Promise<TaskListResponse> {
  const searchParams = new URLSearchParams();

  if (params.limit) {
    searchParams.set("limit", String(params.limit));
  }

  if (params.offset) {
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

  const query = searchParams.toString();

  const endpoint = query
    ? `${taskEndpoints.list()}?${query}`
    : taskEndpoints.list();

  const data = await apiClient<unknown>(endpoint);

  return taskListResponseSchema.parse(data);
}
