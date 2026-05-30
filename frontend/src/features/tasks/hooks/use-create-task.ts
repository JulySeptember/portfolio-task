"use client";
import { toast } from "sonner";

import { useMutation, useQueryClient } from "@tanstack/react-query";

import { createTask, type CreateTaskInput } from "../api/create-task";

import { taskQueryKeys } from "../queries/task-query-keys";

export function useCreateTask() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: CreateTaskInput) => createTask(input),

    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: taskQueryKeys.lists(),
      });
      toast.success("Task Created");
    },

    onError: () => {
      toast.error("Failed to create task");
    },
  });
}
