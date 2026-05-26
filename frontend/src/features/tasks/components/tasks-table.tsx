"use client";

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

import { Button } from "@/components/ui/button";

import { TaskStatusBadge } from "./task-status-badge";

import { useDeleteTask } from "../hooks/use-delete-task";

import { useUpdateTaskStatus } from "../hooks/use-update-task-status";

import type { Task, TaskStatus } from "../schemas/task-schema";
import { EditTaskDialog } from "./edit-task-dialog";

type Props = {
  tasks: Task[];
};

const nextStatusMap: Record<TaskStatus, TaskStatus> = {
  TODO: "DOING",

  DOING: "DONE",

  DONE: "TODO",
};

export function TasksTable({ tasks }: Props) {
  const deleteTask = useDeleteTask();

  const updateStatus = useUpdateTaskStatus();

  if (tasks.length === 0) {
    return (
      <div className="rounded-lg border p-8 text-center">
        <p className="text-muted-foreground text-sm">No tasks found</p>
      </div>
    );
  }

  return (
    <div className="rounded-lg border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Title</TableHead>

            <TableHead>Status</TableHead>

            <TableHead>Due Date</TableHead>

            <TableHead className="w-55">Actions</TableHead>
          </TableRow>
        </TableHeader>

        <TableBody>
          {tasks.map((task) => (
            <TableRow key={task.id}>
              <TableCell>
                <div>
                  <p className="font-medium">{task.title}</p>

                  {task.description && (
                    <p className="text-muted-foreground mt-1 text-sm">
                      {task.description}
                    </p>
                  )}
                </div>
              </TableCell>

              <TableCell>
                <TaskStatusBadge status={task.status} />
              </TableCell>

              <TableCell>
                {task.dueDate ? (
                  <span className="text-sm">
                    {new Date(task.dueDate).toLocaleString()}
                  </span>
                ) : (
                  <span className="text-muted-foreground text-sm">-</span>
                )}
              </TableCell>

              <TableCell>
                <div className="flex items-center gap-2">
                  <Button
                    size="sm"
                    variant="outline"
                    disabled={updateStatus.isPending}
                    onClick={() =>
                      updateStatus.mutate({
                        id: task.id,
                        status: nextStatusMap[task.status],
                      })
                    }
                  >
                    Next Status
                  </Button>

                  <EditTaskDialog task={task} />

                  <Button
                    size="sm"
                    variant="destructive"
                    disabled={deleteTask.isPending}
                    onClick={() => deleteTask.mutate(task.id)}
                  >
                    Delete
                  </Button>
                </div>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}
