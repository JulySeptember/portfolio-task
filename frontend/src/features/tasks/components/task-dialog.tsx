// src/features/tasks/components/task-dialog.tsx
"use client";

import { VisuallyHidden } from "@radix-ui/react-visually-hidden";

import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

import { TaskEditor } from "./task-editor";

import type { Task } from "../schemas/task-schema";

type Props =
  | {
      mode: "create";
      open: boolean;
      onOpenChange: (open: boolean) => void;
      onSuccess?: () => void;
    }
  | {
      mode: "edit";
      task: Task;
      open: boolean;
      onOpenChange: (open: boolean) => void;
      onSuccess?: () => void;
    };

export function TaskDialog(props: Props) {
  return (
    <Dialog
      open={props.open}
      onOpenChange={(nextOpen) => {
        props.onOpenChange(nextOpen);
      }}
    >
      <DialogContent
        onPointerDownOutside={(e) => e.preventDefault()}
        onEscapeKeyDown={(e) => e.preventDefault()}
        className="
          w-[98vw]
          max-w-6xl

          max-h-[80vh]
          md:max-h-[90vh]

          overflow-y-auto
          rounded-2xl

          p-4
          md:p-8                   
          "
      >
        <DialogHeader className="sr-only">
          <VisuallyHidden>
            <DialogTitle>
              {props.mode === "create" ? "Create Task" : "Edit Task"}
            </DialogTitle>
          </VisuallyHidden>
        </DialogHeader>
        {props.mode === "create" ? (
          <TaskEditor
            mode="create"
            onSuccess={() => {
              props.onSuccess?.();
              props.onOpenChange(false);
            }}
          />
        ) : (
          <TaskEditor
            mode="edit"
            task={props.task}
            onSuccess={() => {
              props.onSuccess?.();
              props.onOpenChange(false);
            }}
          />
        )}
      </DialogContent>
    </Dialog>
  );
}
