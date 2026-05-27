"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";

import { toast } from "sonner";

import { deleteTask } from "../api/delete-task";

import { taskQueryKeys } from "../queries/task-query-keys";

import type { Task, TaskListResponse } from "../schemas/task-schema";

export function useDeleteTask() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: deleteTask,

    onMutate: async (id: number) => {
      await queryClient.cancelQueries({
        queryKey: taskQueryKeys.lists(),
      });

      await queryClient.cancelQueries({
        queryKey: taskQueryKeys.detail(id),
      });

      const previousQueries = queryClient.getQueriesData<TaskListResponse>({
        queryKey: taskQueryKeys.lists(),
      });

      const previousTask = queryClient.getQueryData<Task>(
        taskQueryKeys.detail(id),
      );

      queryClient.setQueriesData<TaskListResponse>(
        {
          queryKey: taskQueryKeys.lists(),
        },
        (old) => {
          if (!old) {
            return old;
          }

          const nextItems = old.items.filter((task) => task.id !== id);

          return {
            ...old,

            count: Math.max(
              0,
              old.count - (old.items.length - nextItems.length),
            ),

            items: nextItems,
          };
        },
      );

      queryClient.removeQueries({
        queryKey: taskQueryKeys.detail(id),
      });

      return {
        previousQueries,

        previousTask,
      };
    },

    onError: (error, id, context) => {
      context?.previousQueries.forEach(([queryKey, data]) => {
        queryClient.setQueryData(queryKey, data);
      });

      if (context?.previousTask) {
        queryClient.setQueryData(
          taskQueryKeys.detail(id),
          context.previousTask,
        );
      }

      toast.error(error.message);
    },

    onSuccess: () => {
      toast.success("Task deleted");
    },

    onSettled: async (_, __, id) => {
      await queryClient.invalidateQueries({
        queryKey: taskQueryKeys.lists(),
      });

      await queryClient.invalidateQueries({
        queryKey: taskQueryKeys.detail(id),
      });
    },
  });
}
