// src/features/tasks/hooks/use-update-task.ts

"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";

import { toast } from "sonner";

import { updateTask } from "../api/update-task";

import { taskQueryKeys } from "../queries/task-query-keys";

import type { Task, TaskListResponse } from "../schemas/task-schema";

type UpdateTaskInput = {
  id: number;

  title: string;

  description: string;

  status: "TODO" | "DOING" | "DONE";

  due_date: string | null;
};

type UpdateTaskContext = {
  previousLists: Array<[readonly unknown[], TaskListResponse | undefined]>;
};

export function useUpdateTask() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (input: UpdateTaskInput) => {
      const { id, ...body } = input;

      return updateTask(id, body);
    },

    onMutate: async (
      updatedTask: UpdateTaskInput,
    ): Promise<UpdateTaskContext> => {
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

        const params = queryKey[2] as
          | {
              status?: "TODO" | "DOING" | "DONE";
            }
          | undefined;

        let items = data.items.map((task) =>
          task.id === updatedTask.id
            ? {
                ...task,
                ...updatedTask,
                dueDate: updatedTask.due_date,
              }
            : task,
        );

        if (params?.status) {
          items = items.filter((task) => task.status === params.status);
        }

        queryClient.setQueryData<TaskListResponse>(queryKey, {
          ...data,

          count: items.length,

          items,
        });
      });

      return {
        previousLists,
      };
    },

    onError: (_, __, context?: UpdateTaskContext) => {
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
