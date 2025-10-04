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
import { Badge } from "@/components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useToast } from "@/hooks/use-toast";
import { CheckCircle, XCircle } from "lucide-react";
import { format } from "date-fns";

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
      const api = createClientApprovalAPI(session.accessToken);
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
      const api = createClientApprovalAPI(session.accessToken);
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
          <p className="text-muted-foreground">Review and approve expense requests</p>
        </div>
      </div>

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Employee</TableHead>
              <TableHead>Date</TableHead>
              <TableHead>Category</TableHead>
              <TableHead>Description</TableHead>
              <TableHead>Amount</TableHead>
              <TableHead>Status</TableHead>
              <TableHead className="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {!approvals || approvals.length === 0 ? (
              <TableRow>
                <TableCell colSpan={7} className="text-muted-foreground text-center">
                  No pending approvals
                </TableCell>
              </TableRow>
            ) : (
              approvals.map((approval) => (
                <TableRow key={approval.id}>
                  <TableCell className="font-medium">
                    {approval.expense?.user?.first_name} {approval.expense?.user?.last_name}
                  </TableCell>
                  <TableCell>
                    {approval.expense?.expense_date &&
                      format(new Date(approval.expense.expense_date), "MMM d, yyyy")}
                  </TableCell>
                  <TableCell className="capitalize">{approval.expense?.category}</TableCell>
                  <TableCell className="max-w-xs truncate">
                    {approval.expense?.description}
                  </TableCell>
                  <TableCell>
                    {approval.expense?.amount.toFixed(2)} {approval.expense?.currency}
                  </TableCell>
                  <TableCell>
                    <Badge variant="secondary">{approval.status}</Badge>
                  </TableCell>
                  <TableCell className="text-right">
                    <div className="flex justify-end gap-2">
                      <Button variant="outline" size="sm" onClick={() => handleApprove(approval)}>
                        <CheckCircle className="mr-2 h-4 w-4" />
                        Approve
                      </Button>
                      <Button variant="outline" size="sm" onClick={() => handleReject(approval)}>
                        <XCircle className="mr-2 h-4 w-4" />
                        Reject
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

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
