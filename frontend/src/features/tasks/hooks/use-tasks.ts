import { useQuery } from "@tanstack/react-query";
import { listTasks, type ListTasksParams } from "../api/list-tasks";
import { taskQueryKeys } from "../queries/task-query-keys";
import type { TaskListResponse } from "../schemas/task-schema";

export function useTasks(params?: ListTasksParams) {
  return useQuery<TaskListResponse>({
    queryKey: taskQueryKeys.list(params),
    queryFn: () => listTasks(params),
  });
}
