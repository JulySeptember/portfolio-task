// src/app/(protected)/tasks/new/page.tsx

import Link from "next/link";

import { ArrowLeft } from "lucide-react";

import { Button } from "@/components/ui/button";

import { TaskEditor } from "@/features/tasks/components/task-editor";

export default function NewTaskPage() {
  return (
    <div className="mx-auto w-full max-w-5xl px-4 py-6 sm:px-6 lg:px-8">
      <div className="space-y-10">
        <div className="space-y-5">
          <Button asChild variant="ghost" className="h-auto px-0 text-sm">
            <Link href="/tasks">
              <ArrowLeft className="mr-2 h-4 w-4" />
              Back to Tasks
            </Link>
          </Button>
        </div>

        <TaskEditor
          mode="create"
          showOpenPageButton={false}
          autoResizeDescription
        />
      </div>
    </div>
  );
}
