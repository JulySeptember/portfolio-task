// src/features/tasks/hooks/use-update-task-status.ts
import { useMutation, useQueryClient } from "@tanstack/react-query";

import { toast } from "sonner";

import {
  updateTaskStatus,
  type UpdateTaskStatusInput,
} from "../api/update-task-status";

import { taskQueryKeys } from "../queries/task-query-keys";

import type { Task } from "../schemas/task-schema";

type MutationInput = {
  publicId: string;
  input: UpdateTaskStatusInput;
};

export function useUpdateTaskStatus() {
  const queryClient = useQueryClient();

  return useMutation<Task, Error, MutationInput, { previousTask?: Task }>({
    mutationFn: ({ publicId, input }) => updateTaskStatus(publicId, input),

    onMutate: async ({ publicId, input }) => {
      await queryClient.cancelQueries({
        queryKey: taskQueryKeys.detail(publicId),
      });

      const previousTask = queryClient.getQueryData<Task>(
        taskQueryKeys.detail(publicId),
      );

      if (previousTask) {
        queryClient.setQueryData<Task>(taskQueryKeys.detail(publicId), {
          ...previousTask,
          status: input.status,
        });
      }

      return { previousTask };
    },

    onError: (_error, variables, context) => {
      if (context?.previousTask) {
        queryClient.setQueryData(
          taskQueryKeys.detail(variables.publicId),
          context.previousTask,
        );
      }

      toast.error("Failed to update status");
    },

    onSuccess: (task) => {
      queryClient.setQueryData(taskQueryKeys.detail(task.publicId), task);

      queryClient.invalidateQueries({
        queryKey: taskQueryKeys.lists(),
      });

      toast.success("Task status updated");
    },
  });
}
