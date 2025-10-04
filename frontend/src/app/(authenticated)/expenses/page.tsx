import { auth } from "@/auth";
import { redirect } from "next/navigation";
import ExpensesClient from "@/components/expenses-client";

export default async function Page() {
  const session = await auth();

  if (!session) {
    return redirect("/login");
  }

  return <ExpensesClient />;
}
