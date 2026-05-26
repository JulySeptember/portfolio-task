"use client";

import { Button } from "@/components/ui/button";

type Props = {
  sort: "created_at" | "due_date";

  order: "ASC" | "DESC";

  onSortChange: (sort: "created_at" | "due_date") => void;

  onOrderChange: (order: "ASC" | "DESC") => void;
};

export function TasksSort({ sort, order, onSortChange, onOrderChange }: Props) {
  return (
    <div className="flex flex-wrap items-center gap-2">
      <Button
        variant={sort === "created_at" ? "default" : "outline"}
        onClick={() => onSortChange("created_at")}
      >
        Created
      </Button>

      <Button
        variant={sort === "due_date" ? "default" : "outline"}
        onClick={() => onSortChange("due_date")}
      >
        Due Date
      </Button>

      <Button
        variant="outline"
        onClick={() => onOrderChange(order === "ASC" ? "DESC" : "ASC")}
      >
        {order}
      </Button>
    </div>
  );
}
