// src/features/tasks/hooks/use-update-task.ts

"use client";

import { useMutation } from "@tanstack/react-query";

import { useQueryClient } from "@tanstack/react-query";

import { toast } from "sonner";

import { updateTask } from "../api/update-task";

import { taskQueryKeys } from "../queries/task-query-keys";

import type { Task, TaskListResponse } from "../schemas/task-schema";

export function useUpdateTask() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: updateTask,

    onMutate: async (updatedTask) => {
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
            task.id === updatedTask.id
              ? {
                  ...task,

                  ...updatedTask,

                  dueDate: updatedTask.due_date,
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

      toast.error("Failed to update task");
    },

    onSuccess: (task: Task) => {
      queryClient.setQueryData(taskQueryKeys.detail(task.id), task);

      toast.success("Task updated");
    },

    onSettled: () => {
      queryClient.invalidateQueries({
        queryKey: taskQueryKeys.lists(),
      });
    },
  });
}
