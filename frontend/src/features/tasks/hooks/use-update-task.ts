// src/features/tasks/hooks/use-update-task.ts
"use client";
import { toast } from "sonner";

import { useMutation, useQueryClient } from "@tanstack/react-query";

import {
  updateTask,
  type UpdateTaskInput as ApiUpdateTaskInput,
} from "../api/update-task";

import { taskQueryKeys } from "../queries/task-query-keys";

import type { Task, TaskListResponse } from "../schemas/task-schema";

export type UpdateTaskInput = {
  publicId: string;

  title: string;

  description?: string;

  status: "TODO" | "DOING" | "DONE";

  due_date?: string | null;
};

type UpdateTaskContext = {
  previousLists: Array<[readonly unknown[], TaskListResponse | undefined]>;

  previousDetail?: Task;
};

export function useUpdateTask() {
  const queryClient = useQueryClient();

  return useMutation<Task, Error, UpdateTaskInput, UpdateTaskContext>({
    mutationFn: async (input) => {
      const body: ApiUpdateTaskInput = {
        title: input.title,

        description: input.description ?? "",

        status: input.status,

        due_date: input.due_date ?? null,
      };

      return updateTask(input.publicId, body);
    },

    onMutate: async (updatedTask) => {
      await queryClient.cancelQueries({
        queryKey: taskQueryKeys.lists(),
      });

      await queryClient.cancelQueries({
        queryKey: taskQueryKeys.detail(updatedTask.publicId),
      });

      const previousLists = queryClient.getQueriesData<TaskListResponse>({
        queryKey: taskQueryKeys.lists(),
      });

      const previousDetail = queryClient.getQueryData<Task>(
        taskQueryKeys.detail(updatedTask.publicId),
      );

      // optimistic update (lists)
      previousLists.forEach(([queryKey, data]) => {
        if (!data) {
          return;
        }

        const items = data.items.map((task) =>
          task.publicId === updatedTask.publicId
            ? {
                ...task,

                title: updatedTask.title,

                description: updatedTask.description ?? "",

                status: updatedTask.status,

                dueDate: updatedTask.due_date ?? null,
              }
            : task,
        );

        queryClient.setQueryData<TaskListResponse>(queryKey, {
          ...data,

          items,
        });
      });

      // optimistic update (detail)
      if (previousDetail) {
        queryClient.setQueryData<Task>(
          taskQueryKeys.detail(updatedTask.publicId),
          {
            ...previousDetail,

            title: updatedTask.title,

            description: updatedTask.description ?? "",

            status: updatedTask.status,

            dueDate: updatedTask.due_date ?? null,
          },
        );
      }

      return {
        previousLists,

        previousDetail,
      };
    },

    onError: (_, variables, context) => {
      // rollback lists
      context?.previousLists.forEach(([queryKey, data]) => {
        queryClient.setQueryData(queryKey, data);
      });

      // rollback detail
      if (context?.previousDetail) {
        queryClient.setQueryData(
          taskQueryKeys.detail(variables.publicId),
          context.previousDetail,
        );
      }
      toast.error("Failed to Update task");
    },

    onSuccess: (task) => {
      // detail cache
      queryClient.setQueryData(taskQueryKeys.detail(task.publicId), task);

      // refresh list
      queryClient.invalidateQueries({
        queryKey: taskQueryKeys.lists(),
      });

      toast.success("Task updated");
    },
  });
}
