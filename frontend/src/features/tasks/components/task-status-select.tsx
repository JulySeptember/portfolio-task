// src/features/tasks/components/task-status-select.tsx

"use client";

import { cn } from "@/lib/utils";

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

import type { TaskStatus } from "../schemas/task-schema";

type Props = {
  value: TaskStatus;

  onChange: (value: TaskStatus) => void;

  className?: string;

  disabled?: boolean;
};

function getStatusColor(status: TaskStatus) {
  switch (status) {
    case "TODO":
      return "border-zinc-700 bg-zinc-900 text-zinc-200";

    case "DOING":
      return "border-blue-500/30 bg-blue-500/10 text-blue-300";

    case "DONE":
      return "border-emerald-500/30 bg-emerald-500/10 text-emerald-300";

    default:
      return "";
  }
}

export function TaskStatusSelect({
  value,
  onChange,
  className,
  disabled,
}: Props) {
  return (
    <Select
      value={value}
      onValueChange={(v) => onChange(v as TaskStatus)}
      disabled={disabled}
    >
      <SelectTrigger
        className={cn(
          `
      h-12!
      w-full
      rounded-xl
      px-4

      text-sm
      font-medium

      [&>span]:text-sm!
      [&>span]:font-medium

      transition-colors
    `,
          getStatusColor(value),
          className,
        )}
      >
        <SelectValue />
      </SelectTrigger>

      <SelectContent className="rounded-xl">
        <SelectItem value="TODO" className="h-12 text-sm! font-medium">
          TODO
        </SelectItem>
        <SelectItem
          value="DOING"
          className="h-12 text-sm! font-medium text-blue-300"
        >
          DOING
        </SelectItem>
        <SelectItem
          value="DONE"
          className="h-12 text-sm! font-medium text-emerald-300"
        >
          DONE
        </SelectItem>
      </SelectContent>
    </Select>
  );
}
