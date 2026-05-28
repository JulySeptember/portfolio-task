"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";

import { toast } from "sonner";

import { updateTaskStatus } from "../api/update-task-status";

import { taskQueryKeys } from "../queries/task-query-keys";

import type { TaskListResponse } from "../schemas/task-schema";

type UpdateTaskStatusInput = {
  id: number;

  status: "TODO" | "DOING" | "DONE";
};

export function useUpdateTaskStatus() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, status }: UpdateTaskStatusInput) =>
      updateTaskStatus(id, {
        status,
      }),

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

        const items = data.items.map((task) =>
          task.id === id
            ? {
                ...task,
                status,
              }
            : task,
        );

        queryClient.setQueryData<TaskListResponse>(queryKey, {
          ...data,
          items,
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

      toast.error("Failed to update task status");
    },

    onSuccess: () => {
      toast.success("Task status updated");
    },

    onSettled: () => {
      queryClient.invalidateQueries({
        queryKey: taskQueryKeys.lists(),
      });
    },
  });
}
