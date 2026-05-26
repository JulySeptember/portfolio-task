"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";

import { updateTaskStatus } from "../api/update-task-status";

type Input = {
  id: number;

  status: Parameters<typeof updateTaskStatus>[1];
};

export function useUpdateTaskStatus() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, status }: Input) => updateTaskStatus(id, status),

    onSuccess: async () => {
      await queryClient.invalidateQueries({
        queryKey: ["tasks"],
      });
    },
  });
}
