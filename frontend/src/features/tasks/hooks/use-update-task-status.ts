// src/features/tasks/hooks/use-update-task-status.ts

"use client";

import { useMutation } from "@tanstack/react-query";

import { useQueryClient } from "@tanstack/react-query";

import { toast } from "sonner";

import { updateTaskStatus } from "../api/update-task-status";

import { taskQueryKeys } from "../queries/task-query-keys";

import type { TaskListResponse } from "../schemas/task-schema";

export function useUpdateTaskStatus() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: updateTaskStatus,

    onMutate: async ({ id, status }) => {
      await queryClient.cancelQueries({
        queryKey: taskQueryKeys.lists(),
      });

      const previousLists = queryClient.getQueriesData<TaskListResponse>({
        queryKey: taskQueryKeys.lists(),
      });

      previousLists.forEach(([queryKey, data]) => {
        if (!data) {
          return;
        }

        queryClient.setQueryData<TaskListResponse>(queryKey, {
          ...data,

          items: data.items.map((task) =>
            task.id === id
              ? {
                  ...task,
                  status,
                }
              : task,
          ),
        });
      });

      return {
        previousLists,
      };
    },

    onError: (_, __, context) => {
      context?.previousLists.forEach(([queryKey, data]) => {
        queryClient.setQueryData(queryKey, data);
      });

      toast.error("Failed to update status");
    },

    onSuccess: () => {
      toast.success("Status updated");
    },

    onSettled: () => {
      queryClient.invalidateQueries({
        queryKey: taskQueryKeys.lists(),
      });
    },
  });
}
