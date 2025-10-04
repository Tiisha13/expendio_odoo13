"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import { Approval, ApprovalActionInput } from "@/types/api";
import { createClientApprovalAPI } from "@/lib/api";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { useToast } from "@/hooks/use-toast";
import { DataTable } from "@/components/ui/data-table";
import { createApprovalColumns } from "@/components/approvals-columns";

export default function ApprovalsPage() {
  const { data: session } = useSession();
  const { toast } = useToast();
  const [approvals, setApprovals] = useState<Approval[]>([]);
  const [loading, setLoading] = useState(true);
  const [actionDialogOpen, setActionDialogOpen] = useState(false);
  const [selectedApproval, setSelectedApproval] = useState<Approval | null>(null);
  const [actionType, setActionType] = useState<"approve" | "reject">("approve");
  const [comment, setComment] = useState("");

  const fetchApprovals = async () => {
    if (!session?.accessToken) return;

    try {
      setLoading(true);
      const api = createClientApprovalAPI(session.accessToken, session.user?.refresh_token);
      const response = await api.pending();
      setApprovals(response.data || []); // Ensure it's always an array
    } catch (error: any) {
      setApprovals([]); // Reset to empty array on error
      toast({
        title: "Error",
        description: error.message || "Failed to fetch approvals",
        variant: "destructive",
      });
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchApprovals();
  }, [session]);

  const handleApprove = (approval: Approval) => {
    setSelectedApproval(approval);
    setActionType("approve");
    setComment("");
    setActionDialogOpen(true);
  };

  const handleReject = (approval: Approval) => {
    setSelectedApproval(approval);
    setActionType("reject");
    setComment("");
    setActionDialogOpen(true);
  };

  const handleSubmitAction = async () => {
    if (!session?.accessToken || !selectedApproval) return;

    try {
      const api = createClientApprovalAPI(session.accessToken, session.user?.refresh_token);
      const data: ApprovalActionInput = { comments: comment };

      if (actionType === "approve") {
        await api.approve(selectedApproval.id, data);
      } else {
        await api.reject(selectedApproval.id, data);
      }

      toast({
        title: "Success",
        description: `Expense ${actionType === "approve" ? "approved" : "rejected"} successfully`,
      });

      setActionDialogOpen(false);
      setSelectedApproval(null);
      setComment("");
      fetchApprovals();
    } catch (error: any) {
      toast({
        title: "Error",
        description: error.message || `Failed to ${actionType} expense`,
        variant: "destructive",
      });
    }
  };

  if (loading) {
    return <div className="p-4">Loading...</div>;
  }

  return (
    <div className="space-y-6 p-4">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Pending Approvals</h1>
          <p className="text-muted-foreground">
            Review and approve expense requests ({approvals.length} pending)
          </p>
        </div>
      </div>

      {approvals.length === 0 && !loading && (
        <div className="rounded-lg border border-dashed p-8 text-center">
          <p className="text-muted-foreground">
            No pending approvals assigned to you at this time.
          </p>
        </div>
      )}

      <DataTable
        columns={createApprovalColumns(handleApprove, handleReject)}
        data={approvals}
        searchKey="description"
        searchPlaceholder="Search approvals..."
      />

      <Dialog open={actionDialogOpen} onOpenChange={setActionDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{actionType === "approve" ? "Approve" : "Reject"} Expense</DialogTitle>
            <DialogDescription>
              {selectedApproval && selectedApproval.expense && (
                <div className="mt-2 space-y-2">
                  <p>
                    <strong>Employee:</strong> {selectedApproval.expense.user?.first_name}{" "}
                    {selectedApproval.expense.user?.last_name}
                  </p>
                  <p>
                    <strong>Amount:</strong> {selectedApproval.expense.amount.toFixed(2)}{" "}
                    {selectedApproval.expense.currency}
                  </p>
                  <p>
                    <strong>Description:</strong> {selectedApproval.expense.description}
                  </p>
                </div>
              )}
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="comment">Comment (optional)</Label>
              <Textarea
                id="comment"
                value={comment}
                onChange={(e) => setComment(e.target.value)}
                placeholder="Add a comment..."
                rows={3}
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setActionDialogOpen(false)}>
              Cancel
            </Button>
            <Button
              onClick={handleSubmitAction}
              variant={actionType === "approve" ? "default" : "destructive"}
            >
              {actionType === "approve" ? "Approve" : "Reject"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
