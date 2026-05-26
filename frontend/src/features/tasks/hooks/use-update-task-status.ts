"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";

import { toast } from "sonner";

import { updateTaskStatus } from "../api/update-task-status";

import { taskQueryKeys } from "../queries/task-query-keys";

import type { TaskListResponse } from "../schemas/task-schema";

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
      toast.success("Task status updated");
    },

    onSettled: async () => {
      await queryClient.invalidateQueries({
        queryKey: taskQueryKeys.list(),
      });
    },
  });
}
