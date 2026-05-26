"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";

import { createTask } from "../api/create-task";

import { taskQueryKeys } from "../queries/task-query-keys";

export function useCreateTask() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: createTask,

    onSuccess: async () => {
      await queryClient.invalidateQueries({
        queryKey: taskQueryKeys.list(),
      });
    },
  });
}
