"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";

import { toast } from "sonner";

import { createTask } from "../api/create-task";

import { taskQueryKeys } from "../queries/task-query-keys";

export function useCreateTask() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: createTask,

    onError: (error) => {
      toast.error(error.message);
    },

    onSuccess: async (task) => {
      queryClient.setQueryData(taskQueryKeys.detail(task.id), task);

      toast.success("Task created");

      await queryClient.invalidateQueries({
        queryKey: taskQueryKeys.lists(),
      });
    },
  });
}
