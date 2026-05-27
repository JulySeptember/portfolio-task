"use client";

import { useQuery } from "@tanstack/react-query";

import { getTask } from "../api/get-task";

import { taskQueryKeys } from "../queries/task-query-keys";

export function useTask(id: number, enabled = true) {
  return useQuery({
    queryKey: taskQueryKeys.detail(id),

    queryFn: () => getTask(id),

    enabled: enabled && !!id,
  });
}
