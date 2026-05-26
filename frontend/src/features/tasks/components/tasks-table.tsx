"use client";

import { useState } from "react";

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

import { TaskDetailDialog } from "./task-detail-dialog";

import { useUpdateTaskStatus } from "../hooks/use-update-task-status";

import type { Task, TaskStatus } from "../schemas/task-schema";

type Props = {
  tasks: Task[];
};

export function TasksTable({ tasks }: Props) {
  const [selectedTask, setSelectedTask] = useState<Task | null>(null);

  const updateStatus = useUpdateTaskStatus();

  if (tasks.length === 0) {
    return (
      <div className="rounded-lg border p-8 text-center">
        <p className="text-muted-foreground text-sm">No tasks found</p>
      </div>
    );
  }

  return (
    <>
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
            {tasks.map((task) => (
              <TableRow
                key={task.id}
                className="cursor-pointer hover:bg-muted/50"
                onClick={() => setSelectedTask(task)}
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
      </div>

      {selectedTask && (
        <TaskDetailDialog
          task={selectedTask}
          open={!!selectedTask}
          onOpenChange={(open: boolean) => {
            if (!open) {
              setSelectedTask(null);
            }
          }}
        />
      )}
    </>
  );
}
