// src/features/tasks/server/get-tasks.ts

import { cookies } from "next/headers";

import {
  taskListResponseSchema,
  type TaskListResponse,
} from "../schemas/task-schema";

const API_URL = process.env.NEXT_PUBLIC_API_URL!;

type Params = {
  limit?: number;

  offset?: number;

  status?: "TODO" | "DOING" | "DONE";

  sort?: "created_at" | "due_date";

  order?: "ASC" | "DESC";
};

export async function getTasks(params?: Params): Promise<TaskListResponse> {
  const cookieStore = await cookies();

  const accessToken = cookieStore.get("access_token")?.value;

  if (!accessToken) {
    throw new Error("Unauthorized");
  }

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

  const response = await fetch(`${API_URL}/tasks${query ? `?${query}` : ""}`, {
    method: "GET",

    headers: {
      Authorization: `Bearer ${accessToken}`,
    },

    cache: "no-store",
  });

  if (!response.ok) {
    throw new Error("Failed to fetch tasks");
  }

  const data = await response.json();

  return taskListResponseSchema.parse(data);
}
