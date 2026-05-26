"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";

import { updateTask } from "../api/update-task";

type Input = {
  id: number;

  data: Parameters<typeof updateTask>[1];
};

export function useUpdateTask() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: Input) => updateTask(id, data),

    onSuccess: async () => {
      await queryClient.invalidateQueries({
        queryKey: ["tasks"],
      });
    },
  });
}
