"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";

import { toast } from "sonner";

import { createTask } from "../api/create-task";

import { taskQueryKeys } from "../queries/task-query-keys";

import type {
  Task,
  TaskListResponse,
  TaskRequest,
  TaskStatus,
} from "../schemas/task-schema";

type ListParams = {
  status?: TaskStatus;

  sort?: "created_at" | "due_date";

  order?: "ASC" | "DESC";

  limit?: number;

  offset?: number;
};

export function useCreateTask() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: createTask,

    onMutate: async (newTask: TaskRequest) => {
      await queryClient.cancelQueries({
        queryKey: taskQueryKeys.lists(),
      });

      const previousQueries = queryClient.getQueriesData<TaskListResponse>({
        queryKey: taskQueryKeys.lists(),
      });

      const optimisticTask: Task = {
        id: Date.now(),

        userId: 0,

        title: newTask.title,

        description: newTask.description,

        status: newTask.status,

        dueDate: newTask.due_date,

        createdAt: new Date().toISOString(),

        updatedAt: new Date().toISOString(),
      };

      previousQueries.forEach(([queryKey, data]) => {
        if (!data) {
          return;
        }

        const params = queryKey[2] as ListParams | undefined;

        // filter mismatch
        if (params?.status && params.status !== newTask.status) {
          return;
        }

        queryClient.setQueryData<TaskListResponse>(queryKey, {
          ...data,

          count: data.count + 1,

          items: [optimisticTask, ...data.items].slice(0, data.limit),
        });
      });

      return {
        previousQueries,
      };
    },

    onError: (error, _, context) => {
      context?.previousQueries.forEach(([queryKey, data]) => {
        queryClient.setQueryData(queryKey, data);
      });

      toast.error(error.message);
    },

    onSuccess: async (task) => {
      queryClient.setQueryData(taskQueryKeys.detail(task.id), task);

      toast.success("Task created");

      await queryClient.invalidateQueries({
        queryKey: taskQueryKeys.lists(),
      });
    },

    onSettled: async () => {
      await queryClient.invalidateQueries({
        queryKey: taskQueryKeys.lists(),
      });
    },
  });
}
