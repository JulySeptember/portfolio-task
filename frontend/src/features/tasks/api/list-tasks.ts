import { apiClient } from "@/lib/api/client";

import {
  taskListSchema,
  type TaskListResponse,
  type TaskStatus,
} from "../schemas/task-schema";

// =========================
// params
// =========================

export type ListTasksParams = {
  status?: TaskStatus;

  sort?: "created_at" | "due_date";

  order?: "ASC" | "DESC";

  limit?: number;

  offset?: number;
};

// =========================
// list tasks
// =========================

export async function listTasks(
  params?: ListTasksParams,
): Promise<TaskListResponse> {
  const searchParams = new URLSearchParams();

  if (params?.status) {
    searchParams.set("status", params.status);
  }

  if (params?.sort) {
    searchParams.set("sort", params.sort);
  }

  if (params?.order) {
    searchParams.set("order", params.order);
  }

  if (typeof params?.limit === "number") {
    searchParams.set("limit", String(params.limit));
  }

  if (typeof params?.offset === "number") {
    searchParams.set("offset", String(params.offset));
  }

  const query = searchParams.toString();

  const url = query
    ? `${process.env.NEXT_PUBLIC_API_URL}/tasks?${query}`
    : `${process.env.NEXT_PUBLIC_API_URL}/tasks`;

  const data = await apiClient<unknown>(url);

  return taskListSchema.parse(data);
}
