"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";

import { toast } from "sonner";

import { updateTask } from "../api/update-task";

import { taskQueryKeys } from "../queries/task-query-keys";

import type {
  Task,
  TaskListResponse,
  TaskRequest,
} from "../schemas/task-schema";

type Variables = {
  id: number;
} & TaskRequest;

export function useUpdateTask() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: updateTask,

    onMutate: async (updatedTask: Variables) => {
      await queryClient.cancelQueries({
        queryKey: taskQueryKeys.lists(),
      });

      await queryClient.cancelQueries({
        queryKey: taskQueryKeys.detail(updatedTask.id),
      });

      const previousQueries = queryClient.getQueriesData<TaskListResponse>({
        queryKey: taskQueryKeys.lists(),
      });

      const previousTask = queryClient.getQueryData<Task>(
        taskQueryKeys.detail(updatedTask.id),
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
              task.id === updatedTask.id
                ? {
                    ...task,

                    title: updatedTask.title,

                    description: updatedTask.description,

                    status: updatedTask.status,

                    dueDate: updatedTask.due_date,

                    updatedAt: new Date().toISOString(),
                  }
                : task,
            ),
          };
        },
      );

      queryClient.setQueryData(
        taskQueryKeys.detail(updatedTask.id),
        (old: Task | undefined) => {
          if (!old) {
            return old;
          }

          return {
            ...old,

            title: updatedTask.title,

            description: updatedTask.description,

            status: updatedTask.status,

            dueDate: updatedTask.due_date,

            updatedAt: new Date().toISOString(),
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
      toast.success("Task updated");
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
