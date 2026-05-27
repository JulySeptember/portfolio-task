// src/features/tasks/server/get-task.ts

import { cookies } from "next/headers";

import { taskSchema, type Task } from "../schemas/task-schema";

const API_URL = process.env.NEXT_PUBLIC_API_URL!;

export async function getTask(id: number): Promise<Task> {
  const cookieStore = await cookies();

  const accessToken = cookieStore.get("access_token")?.value;

  if (!accessToken) {
    throw new Error("Unauthorized");
  }

  const response = await fetch(`${API_URL}/tasks/${id}`, {
    method: "GET",

    headers: {
      Authorization: `Bearer ${accessToken}`,
    },

    cache: "no-store",
  });

  if (!response.ok) {
    throw new Error("Failed to fetch task");
  }

  const data = await response.json();

  return taskSchema.parse(data);
}
