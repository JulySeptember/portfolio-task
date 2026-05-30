"use client";

import { toast } from "sonner";

import { useMutation, useQueryClient } from "@tanstack/react-query";

import { deleteTask } from "../api/delete-task";

import { taskQueryKeys } from "../queries/task-query-keys";

export function useDeleteTask() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (publicId: string) => deleteTask(publicId),

    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: taskQueryKeys.lists(),
      });

      toast.success("Task Deleted");
    },
    onError: () => {
      toast.error("Failed to Delete task");
    },
  });
}
