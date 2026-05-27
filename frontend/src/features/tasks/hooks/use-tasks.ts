"use client";

import { useQuery } from "@tanstack/react-query";

import { listTasks, type ListTasksParams } from "../api/list-tasks";

import { taskQueryKeys } from "../queries/task-query-keys";

export function useTasks(params?: ListTasksParams) {
  return useQuery({
    queryKey: taskQueryKeys.list(params),

    queryFn: () => listTasks(params),
  });
}
