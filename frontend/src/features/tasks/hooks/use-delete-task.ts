"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";

import { deleteTask } from "../api/delete-task";

export function useDeleteTask() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: deleteTask,

    onSuccess: async () => {
      await queryClient.invalidateQueries({
        queryKey: ["tasks"],
      });
    },
  });
}
