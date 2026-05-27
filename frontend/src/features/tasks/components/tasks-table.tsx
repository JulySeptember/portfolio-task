// src/features/tasks/components/tasks-table.tsx

"use client";

import { useRouter, useSearchParams } from "next/navigation";

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

import { useUpdateTaskStatus } from "../hooks/use-update-task-status";

import { useTasks } from "../hooks/use-tasks";

import type {
  Task,
  TaskListResponse,
  TaskStatus,
} from "../schemas/task-schema";

import { TasksPagination } from "./tasks-pagination";

type Props = {
  initialData: TaskListResponse;

  limit: number;

  offset: number;

  status?: "TODO" | "DOING" | "DONE";

  sort?: "created_at" | "due_date";

  order?: "ASC" | "DESC";
};

export function TasksTable({
  initialData,
  limit,
  offset,
  status,
  sort,
  order,
}: Props) {
  const router = useRouter();

  const searchParams = useSearchParams();

  const updateStatus = useUpdateTaskStatus();

  const { data } = useTasks(
    {
      limit,
      offset,
      status,
      sort,
      order,
    },
    {
      initialData,
    },
  );

  const tasks: Task[] = data?.items ?? [];

  if (tasks.length === 0) {
    return (
      <div className="rounded-xl border p-8 text-center">
        <p className="text-muted-foreground text-sm">No tasks found</p>

        <div className="border-t px-6 pt-6">
          <TasksPagination
            total={data?.count ?? 0}
            limit={data?.limit ?? limit}
            offset={data?.offset ?? offset}
          />
        </div>
      </div>
    );
  }

  return (
    <div className="rounded-xl border">
      <Table className="table-fixed">
        <TableHeader>
          <TableRow>
            <TableHead className="w-[55%]">Title</TableHead>

            <TableHead className="w-45">Status</TableHead>

            <TableHead className="w-55">Due Date</TableHead>
          </TableRow>
        </TableHeader>

        <TableBody>
          {tasks.map((task: Task) => (
            <TableRow
              key={task.id}
              className="cursor-pointer hover:bg-muted/50"
              onClick={() => {
                const params = new URLSearchParams(searchParams.toString());

                params.set("taskId", String(task.id));

                router.push(`/tasks?${params.toString()}`);
              }}
            >
              <TableCell>
                <div className="min-w-0 space-y-1">
                  <p className="truncate font-medium">{task.title}</p>

                  {task.description && (
                    <p className="text-muted-foreground line-clamp-2 break-all text-sm">
                      {task.description}
                    </p>
                  )}
                </div>
              </TableCell>

              <TableCell onClick={(e) => e.stopPropagation()}>
                <Select
                  value={task.status}
                  onValueChange={(value: TaskStatus) => {
                    updateStatus.mutate({
                      id: task.id,

                      status: value,
                    });
                  }}
                >
                  <SelectTrigger className="w-32">
                    <SelectValue />
                  </SelectTrigger>

                  <SelectContent>
                    <SelectItem value="TODO">TODO</SelectItem>

                    <SelectItem value="DOING">DOING</SelectItem>

                    <SelectItem value="DONE">DONE</SelectItem>
                  </SelectContent>
                </Select>
              </TableCell>

              <TableCell>
                {task.dueDate ? (
                  <span className="truncate text-sm">
                    {new Date(task.dueDate).toLocaleString()}
                  </span>
                ) : (
                  <span className="text-muted-foreground text-sm">-</span>
                )}
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>

      <div className="border-t px-6 py-4">
        <TasksPagination
          total={data?.count ?? 0}
          limit={data?.limit ?? limit}
          offset={data?.offset ?? offset}
        />
      </div>
    </div>
  );
}
