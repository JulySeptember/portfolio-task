// src/app/(protected)/tasks/[id]/page.tsx

import Link from "next/link";

import { notFound } from "next/navigation";

import { ArrowLeft } from "lucide-react";

import { Button } from "@/components/ui/button";

import { getTask } from "@/features/tasks/server/get-task";

import { TaskEditor } from "@/features/tasks/components/task-editor";

type Props = {
  params: Promise<{
    id: string;
  }>;
};

export default async function TaskDetailPage({ params }: Props) {
  const { id } = await params;

  const taskId = Number(id);

  if (!Number.isInteger(taskId) || taskId <= 0) {
    notFound();
  }

  const task = await getTask(taskId);

  return (
    <div className="mx-auto w-full max-w-5xl px-4 py-6 sm:px-6 lg:px-8">
      <div className="space-y-10">
        {/* top */}

        <div className="space-y-5">
          <Button asChild variant="ghost" className="h-auto px-0 text-sm">
            <Link href="/tasks">
              <ArrowLeft className="mr-2 h-4 w-4" />
              Back to Tasks
            </Link>
          </Button>
        </div>

        {/* editor */}

        <TaskEditor mode="edit" task={task} autoResizeDescription />
      </div>
    </div>
  );
}
