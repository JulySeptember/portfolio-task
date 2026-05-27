"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";

import { toast } from "sonner";

import { updateTaskStatus } from "../api/update-task-status";

import { taskQueryKeys } from "../queries/task-query-keys";

import type { Task, TaskListResponse } from "../schemas/task-schema";

type Variables = {
  id: number;

  status: "TODO" | "DOING" | "DONE";
};

export function useUpdateTaskStatus() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: updateTaskStatus,

    onMutate: async ({ id, status }: Variables) => {
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

          return {
            ...old,

            items: old.items.map((task) =>
              task.id === id
                ? {
                    ...task,

                    status,
                  }
                : task,
            ),
          };
        },
      );

      queryClient.setQueryData(
        taskQueryKeys.detail(id),
        (old: Task | undefined) => {
          if (!old) {
            return old;
          }

          return {
            ...old,

            status,
          };
        },
      );

      return {
        previousQueries,

        previousTask,
      };
    },

    onError: (error, variables, context) => {
      context?.previousQueries.forEach(([queryKey, data]) => {
        queryClient.setQueryData(queryKey, data);
      });

      if (context?.previousTask) {
        queryClient.setQueryData(
          taskQueryKeys.detail(variables.id),
          context.previousTask,
        );
      }

      toast.error(error.message);
    },

    onSuccess: () => {
      toast.success("Task status updated");
    },

    onSettled: async (_, __, variables) => {
      await queryClient.invalidateQueries({
        queryKey: taskQueryKeys.lists(),
      });

      await queryClient.invalidateQueries({
        queryKey: taskQueryKeys.detail(variables.id),
      });
    },
  });
}
