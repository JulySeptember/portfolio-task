// src/features/tasks/hooks/use-delete-task.ts

"use client";

import { useMutation } from "@tanstack/react-query";

import { useQueryClient } from "@tanstack/react-query";

import { toast } from "sonner";

import { deleteTask } from "../api/delete-task";

import { taskQueryKeys } from "../queries/task-query-keys";

import type { TaskListResponse } from "../schemas/task-schema";

export function useDeleteTask() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: deleteTask,

    onMutate: async (taskId) => {
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

          count: Math.max(0, data.count - 1),

          items: data.items.filter((task) => task.id !== taskId),
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

      toast.error("Failed to delete task");
    },

    onSuccess: () => {
      toast.success("Task deleted");
    },

    onSettled: () => {
      queryClient.invalidateQueries({
        queryKey: taskQueryKeys.lists(),
      });
    },
  });
}
