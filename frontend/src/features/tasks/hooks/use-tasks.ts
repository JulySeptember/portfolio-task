"use client";

import { useQuery } from "@tanstack/react-query";

import { listTasks, type ListTasksParams } from "../api/list-tasks";

export function useTasks(params?: ListTasksParams) {
  return useQuery({
    queryKey: ["tasks", params],

    queryFn: () => listTasks(params),
  });
}
