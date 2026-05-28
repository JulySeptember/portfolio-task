"use client";

import { useQuery } from "@tanstack/react-query";

import { listTasks, type ListTasksParams } from "../api/list-tasks";

import { taskQueryKeys } from "../queries/task-query-keys";

import type { TaskListResponse } from "../schemas/task-schema";

type Options = {
  initialData?: TaskListResponse;
};

export function useTasks(params?: ListTasksParams, options?: Options) {
  return useQuery<TaskListResponse>({
    queryKey: taskQueryKeys.list(params),

    queryFn: () => listTasks(params),

    initialData: options?.initialData,
  });
}
