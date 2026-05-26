"use client";

import { Button } from "@/components/ui/button";

type Props = {
  status: "TODO" | "DOING" | "DONE" | undefined;

  onChange: (status: "TODO" | "DOING" | "DONE" | undefined) => void;
};

export function TasksFilter({ status, onChange }: Props) {
  return (
    <div className="flex flex-wrap items-center gap-2">
      <Button
        variant={!status ? "default" : "outline"}
        onClick={() => onChange(undefined)}
      >
        ALL
      </Button>

      <Button
        variant={status === "TODO" ? "default" : "outline"}
        onClick={() => onChange("TODO")}
      >
        TODO
      </Button>

      <Button
        variant={status === "DOING" ? "default" : "outline"}
        onClick={() => onChange("DOING")}
      >
        DOING
      </Button>

      <Button
        variant={status === "DONE" ? "default" : "outline"}
        onClick={() => onChange("DONE")}
      >
        DONE
      </Button>
    </div>
  );
}
