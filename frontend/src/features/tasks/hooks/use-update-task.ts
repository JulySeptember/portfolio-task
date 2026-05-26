"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";

import { toast } from "sonner";

import { updateTask } from "../api/update-task";

import { taskQueryKeys } from "../queries/task-query-keys";

import type { TaskListResponse, TaskRequest } from "../schemas/task-schema";

type Variables = {
  id: number;
} & TaskRequest;

export function useUpdateTask() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: updateTask,

    onMutate: async (updatedTask: Variables) => {
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
              task.id === updatedTask.id
                ? {
                    ...task,

                    title: updatedTask.title,

                    description: updatedTask.description,

                    status: updatedTask.status,

                    dueDate: updatedTask.due_date,
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
      toast.success("Task updated");
    },

    onSettled: async () => {
      await queryClient.invalidateQueries({
        queryKey: taskQueryKeys.list(),
      });
    },
  });
}
