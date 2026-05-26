"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";

import { toast } from "sonner";

import { createTask } from "../api/create-task";

import { taskQueryKeys } from "../queries/task-query-keys";

export function useCreateTask() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: createTask,

    onSuccess: async () => {
      toast.success("Task created");

      await queryClient.invalidateQueries({
        queryKey: taskQueryKeys.list(),
      });
    },

    onError: (error: Error) => {
      toast.error(error.message);
    },
  });
}
