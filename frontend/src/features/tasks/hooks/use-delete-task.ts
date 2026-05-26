"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";

import { toast } from "sonner";

import { deleteTask } from "../api/delete-task";

import { taskQueryKeys } from "../queries/task-query-keys";

import type { TaskListResponse } from "../schemas/task-schema";

export function useDeleteTask() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: deleteTask,

    onMutate: async (id: number) => {
      await queryClient.cancelQueries({
        queryKey: taskQueryKeys.list(),
      });

      const previousTasks = queryClient.getQueryData<TaskListResponse>(
        taskQueryKeys.list(),
      );

      queryClient.setQueryData<TaskListResponse>(
        taskQueryKeys.list(),
        (old) => {
          if (!old) {
            return old;
          }

          return {
            ...old,

            count: old.count - 1,

            items: old.items.filter((task) => task.id !== id),
          };
        },
      );

      return {
        previousTasks,
      };
    },

    onError: (error, _, context) => {
      if (context?.previousTasks) {
        queryClient.setQueryData(taskQueryKeys.list(), context.previousTasks);
      }

      toast.error(error.message);
    },

    onSuccess: () => {
      toast.success("Task deleted");
    },

    onSettled: async () => {
      await queryClient.invalidateQueries({
        queryKey: taskQueryKeys.list(),
      });
    },
  });
}
